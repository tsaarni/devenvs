# Integrated Storage / Raft Compaction Disk Usage Test

## Overview

This test shows how Raft logstore compaction affects disk usage. In the test the same client workload is run against two server configurations: one with snapshots/compaction disabled and one with snapshots enabled.

Related issue: Improve handling and recovery from disk full condition (see [openbao#1233](https://github.com/openbao/openbao/issues/1233)).

**Test overview**

1. Write 10,000 secrets (fills the Raft log store `raft.db`).
2. Wait 10 seconds for snapshots/compaction to run (if configured).
3. Delete 10,000 secrets.
4. Wait 10 seconds again for snapshots/compaction to run (if configured).
5. Start again from step 1 to repeat the cycle (for monitoring trends over multiple cycles).

During each step, metrics are collected and printed to the console.

**Environment**

- Server: OpenBao 2.5.1 inside Docker. Tests run with a temporary data directory (`--tmpfs /data`, and can be also space constrained by adding e.g. `--tmpfs /data:size=128m`). Each run starts from an empty database.
- Client: A Go program ([`main.go`](main.go)) that sends deterministic traffic and prints metrics. The rate is limited to 500 operations per second to make disk space usage and compaction effects easier to observe.



## Results
For sample output, see [RESULTS.txt](RESULTS.txt), which includes three repeated test cycles for both configurations.

Summary of results:

**Test 1:** No snapshots/compaction during test (`snapshot_threshold=999999`, `snapshot_interval=999999s`):

- `raft.db` grows by ~144 MiB per write+delete cycle and never shrinks.
- After 3 cycles: 16 MiB -> 144 MiB -> 288 MiB -> 432 MiB...
- `RAFT_FREE` stays near zero, no pages are ever freed or reused in the logstore.
- `vault.db` stabilizes at 144 MiB after the first cycle; subsequent cycles reuse its pages.

**Test 2:** Aggressive snapshots/compaction (`snapshot_threshold=100`, `snapshot_interval=5s`, `trailing_logs=10`):

- `raft.db` grows to 64 MiB during the first cycle and then **stabilizes**, pages are continuously compacted and reused over all cycles.
- `RAFT_FREE` shows 13–16 KiB of free pages available for reuse at all times.
- `vault.db` behaves identically to the no-compaction case, stabilizing at 144 MiB.



## Metrics collected

| Metric | Description |
|--------|-------------|
| RAFT_FILE | `/data/raft.db` (logstore) file size on disk |
| VAULT_FILE | `/data/raft/vault.db` (FSM) file size on disk |
| RAFT_FREE | free pages in `/data/raft.db` (available for reuse) |
| RAFT_ALLOC | bytes allocated to `/data/raft.db` freelist structure |
| RAFT_USED | bytes used by `/data/raft.db` freelist metadata |
| RAFT_PEND | pages pending release in `/data/raft.db` |
| FSM_FREE | free pages in `/data/raft/vault.db` (available for reuse) |
| FSM_ALLOC | bytes allocated to `/data/raft/vault.db` freelist structure |
| FSM_USED | bytes used by `/data/raft/vault.db` freelist metadata |
| FSM_PEND | pages pending release in `/data/raft/vault.db` |
| SECRETS | current number of secrets stored |
| OPS | operations (writes+deletes) since last printout |

## Test execution steps

1) Without snapshots/compaction during test
   - Server config: [`raft-no-compaction.hcl`](raft-no-compaction.hcl).
   - Expected behaviour: Raft log grows without freeing pages, `/data/raft.db` keeps increasing.

   Run server example:
   ```bash
   docker run \
     --name openbao-raft-test \
     --tmpfs /data \
     --volume ./:/input:ro \
     --publish 8200:8200 \
     --rm -it \
     ghcr.io/openbao/openbao:2.5.1 \
     bao server -config /input/raft-no-compaction.hcl -log-level=trace
   ```

   Initialize and unseal the server before running the client.
   ```bash
   ./setup.sh
   ```

   The script writes the credentials to `init.json` which the client uses to authenticate.

   Run the test client several times to see the trends in metrics printed to the console:
   ```bash
   go run .
   ```

2) Very agressive snapshots/compaction
   - Server config: [`raft-with-compaction.hcl`](raft-with-compaction.hcl).
   - Expected behaviour: automatic snapshots compact the Raft log, freeing pages for reuse; `/data/raft.db` grows only initially and then stabilizes.

   Run server example:
   ```bash
   docker run \
     --name openbao-raft-test \
     --tmpfs /data \
     --volume ./:/input:ro \
     --publish 8200:8200 \
     --rm -it \
     ghcr.io/openbao/openbao:2.5.1 \
     bao server -config /input/raft-with-compaction.hcl -log-level=trace
   ```

   Initialize and unseal the server before running the client:
   ```bash
   ./setup.sh
   ```

   The script writes the credentials to `init.json` which the client uses to authenticate.

   Run the test client several times to see the trends in metrics printed to the console:

   ```bash
   go run .
   ```

   Observe the server log messages that show the compaction process is happening at the expected intervals:
   ```
   starting snapshot up to: ...
   compacting logs: ...
   snapshot complete up to: ...
   ```
