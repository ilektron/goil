package main

import (
	"fmt"
	"html"
	"log"
	"os"

	"golang.org/x/exp/inotify"
	//	"github.com/pkg/errors"
)

func main() {

	if len(os.Args) > 1 {

	}
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Watch("/tmp")
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case ev := <-watcher.Event:
			log.Println("event:", ev)
		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
}
