package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/tsaarni/echoclient/worker"
)

const (
	baoAddr   = "http://127.0.0.1:8200"
	container = "openbao-raft-test"
)

var (
	rootToken   string
	opCount     atomic.Int64
	storedCount atomic.Int64
	opsSince    atomic.Int64
)

func main() {
	data, err := os.ReadFile("init.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read init.json: %v\n", err)
		os.Exit(1)
	}
	var init_ map[string]any
	json.Unmarshal(data, &init_)
	rootToken = init_["root_token"].(string)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go monitor(ctx)

	fmt.Println("\n--- INITIAL ---")
	printHeader()
	printReport()

	payload := strings.Repeat("x", 13000)

	writeFunc := func(_ context.Context, _ *worker.WorkerPool) error {
		n := opCount.Add(1)
		apiCall("PUT", fmt.Sprintf("/v1/secret/key-%d", n), map[string]any{"value": payload})
		storedCount.Add(1)
		opsSince.Add(1)
		return nil
	}

	deleteFunc := func(_ context.Context, _ *worker.WorkerPool) error {
		n := opCount.Add(1)
		apiCall("DELETE", fmt.Sprintf("/v1/secret/key-%d", n), nil)
		storedCount.Add(-1)
		opsSince.Add(1)
		return nil
	}

	phase := func(msg string) worker.LifeCycleFunc {
		return func(_ context.Context, _ *worker.WorkerPool) {
			opCount.Store(0)
			fmt.Println(msg)
		}
	}

	onEnd := func(_ context.Context, _ *worker.WorkerPool) {
		printReport()
	}

	pool := worker.NewMultiStepWorkerPool(writeFunc, []*worker.Step{
		worker.NewStep(
			worker.WithConcurrency(10),
			worker.WithRepetitions(10000),
			worker.WithRateLimit(500, 500),
			worker.WithHooks(phase("\n--- PHASE 1: Writing 10000 secrets ---"), onEnd),
		),
		worker.NewStep(
			worker.WithPause(10*time.Second),
			worker.WithHooks(func(_ context.Context, _ *worker.WorkerPool) {
				fmt.Println("\n--- PHASE 2: Idling (10 seconds) ---")
			}, onEnd),
		),
		worker.NewStep(
			worker.WithConcurrency(10),
			worker.WithRepetitions(10000),
			worker.WithRateLimit(500, 500),
			worker.WithWorkerFunc(deleteFunc),
			worker.WithHooks(phase("\n--- PHASE 3: Deleting all secrets ---"), onEnd),
		),
		worker.NewStep(
			worker.WithPause(10*time.Second),
			worker.WithHooks(func(_ context.Context, _ *worker.WorkerPool) {
				fmt.Println("\n--- PHASE 4: Idling (10 seconds) ---")
			}, onEnd),
		),
	})
	if err := pool.Launch(); err != nil {
		panic(err)
	}
	pool.Wait()
}

func printHeader() {
	fmt.Printf("%10s %10s %10s %10s %10s %10s %10s %10s %10s %10s %10s %8s %6s\n",
		"TIME", "RAFT_FILE", "VAULT_FILE", "RAFT_FREE", "RAFT_ALLOC", "RAFT_USED", "RAFT_PEND", "FSM_FREE", "FSM_ALLOC", "FSM_USED", "FSM_PEND", "SECRETS", "OPS")
}

func monitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(3 * time.Second):
		}

		printReport()
	}
}

func printReport() {
	newSince := opsSince.Swap(0)
	m := fetchMetrics()
	fmt.Printf("%10s %10s %10s %10s %10s %10s %10s %10s %10s %10s %10s %8d %6d\n",
		time.Now().Format("15:04:05"),
		humanize.IBytes(uint64(containerFileSize("/data/raft/raft.db"))),
		humanize.IBytes(uint64(containerFileSize("/data/vault.db"))),
		humanizeMetric(m["raft_free"]), humanizeMetric(m["raft_alloc"]), humanizeMetric(m["raft_used"]), m["raft_pend"],
		humanizeMetric(m["fsm_free"]), humanizeMetric(m["fsm_alloc"]), humanizeMetric(m["fsm_used"]), m["fsm_pend"],
		storedCount.Load(), newSince,
	)
}

func humanizeMetric(s string) string {
	if s == "?" {
		return s
	}
	if v, err := strconv.ParseUint(s, 10, 64); err == nil {
		return humanize.IBytes(v)
	}
	return s
}

func containerFileSize(path string) int64 {
	out, err := exec.Command("docker", "exec", container, "stat", "-c", "%s", path).Output()
	if err != nil {
		return 0
	}
	n, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	return n
}

var metricRe = regexp.MustCompile(`vault_raft_storage_bolt_freelist_(free_pages|allocated_bytes|used_bytes|pending_pages)\{[^}]*database="(logstore|fsm)"[^}]*\}\s+([0-9.e+]+)`)

func fetchMetrics() map[string]string {
	m := map[string]string{
		"raft_free": "?", "raft_alloc": "?", "raft_used": "?", "raft_pend": "?",
		"fsm_free": "?", "fsm_alloc": "?", "fsm_used": "?", "fsm_pend": "?",
	}

	req, _ := http.NewRequest("GET", baoAddr+"/v1/sys/metrics?format=prometheus", nil)
	req.Header.Set("X-Vault-Token", rootToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return m
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	keys := map[[2]string]string{
		{"free_pages", "logstore"}:      "raft_free",
		{"allocated_bytes", "logstore"}: "raft_alloc",
		{"used_bytes", "logstore"}:      "raft_used",
		{"pending_pages", "logstore"}:   "raft_pend",
		{"free_pages", "fsm"}:           "fsm_free",
		{"allocated_bytes", "fsm"}:      "fsm_alloc",
		{"used_bytes", "fsm"}:           "fsm_used",
		{"pending_pages", "fsm"}:        "fsm_pend",
	}

	for _, match := range metricRe.FindAllStringSubmatch(string(body), -1) {
		if key, ok := keys[[2]string{match[1], match[2]}]; ok {
			if f, err := strconv.ParseFloat(match[3], 64); err == nil {
				m[key] = fmt.Sprintf("%.0f", f)
			}
		}
	}
	return m
}

func apiCall(method, path string, body map[string]any) map[string]any {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, baoAddr+path, reqBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if rootToken != "" {
		req.Header.Set("X-Vault-Token", rootToken)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		return nil
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var result map[string]any
	json.Unmarshal(data, &result)
	return result
}
