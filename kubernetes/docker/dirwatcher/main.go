package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var numSuccess int
var numFailure int

func main() {

	// Check if debug mode is enabled.
	debug := os.Getenv("DEBUG") == "true"

	log.Println("Watching paths: /secret/, /var/run/secrets/kubernetes.io/serviceaccount/")

	// Create separate watchers for secret and projected volumes.  That allows us to react quicker to the changes.
	watcher1, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher1.Close()
	err = watcher1.Add("/secret/")
	if err != nil {
		log.Fatal(err)
	}
	go processWatchEvents(watcher1, "/secret/password", debug)

	watcher2, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher2.Close()
	err = watcher2.Add("/var/run/secrets/kubernetes.io/serviceaccount/")
	if err != nil {
		log.Fatal(err)
	}
	go processWatchEvents(watcher2, "/var/run/secrets/kubernetes.io/serviceaccount/token", debug)

	// Block main goroutine forever.
	<-make(chan struct{})
}

func processWatchEvents(watcher *fsnotify.Watcher, filename string, debug bool) {
	baseDir := filepath.Dir(filename)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if debug {
				log.Println("    Event:", event.Op.String(), event.Name)
			}

			// Wait until files are updated i.e. kubelet calls Rename("..data_tmp", "..data").
			// Fsnotify's Rename equals to IN_MOVED_TO inotify event.

			if event.Has(fsnotify.Rename) {

				// Attempt to open file.
				f, err := os.Open(filename)
				if err != nil {
					numFailure++
					log.Printf("Error : %v\n", err)
				} else {
					f.Close()
					numSuccess++
				}

				log.Printf("Open: succeeded: %d, failed: %d\n", numSuccess, numFailure)

				if debug {
					// List all files in the directory.
					filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
						if err != nil {
							log.Println(err)
							return err
						}
						info, err := os.Stat(path)
						if err != nil {
							log.Println(err)
							return err
						}
						stat := info.Sys().(*syscall.Stat_t)

						log.Printf("    Dirwalk: %s %8d %8d %s\n", info.Mode().String(), stat.Uid, stat.Gid, path)
						return nil
					})
					log.Println("    Note 1: due to race condition the permissions might change during the walk.")
					log.Println("    Note 2: string representation of mode bits are not the same in go and linux https://pkg.go.dev/os#pkg-constants")
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
