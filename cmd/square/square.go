package main

import (
	"log"
	"os"

	"golang.org/x/exp/inotify"
	//	"github.com/pkg/errors"
)

type process func(string)

func main() {
	// Watch all of the arguments passed in

	messages := make(chan string)
	num_watches := len(os.Args - 1)

	for _, arg := range os.Args[1:] {
		log.Println("Arg:", arg)
		go watch_item(arg, messages)
	}



}

func watch_item(path string, fn process, chan done) {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Watch(path)
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

func square_video(video string) {
	log.Println("Processing video", video)
}
