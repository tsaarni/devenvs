package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Rename) {
					// Try read file
					_, err := os.ReadFile(os.Args[2])
					if err != nil {
						log.Printf("CANNOT READ FILE: %v\n", err)
					}

					// List all files in the directory.
					filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
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

						log.Printf("%s %8d %8d %s\n", info.Mode().String(), stat.Uid, stat.Gid, path)
						return nil
					})

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	log.Println("Watching path:", os.Args[1])
	err = watcher.Add(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Testing readability by reading:", os.Args[2])

	// Block main goroutine forever.
	<-make(chan struct{})
}
