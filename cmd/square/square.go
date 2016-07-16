package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"

	"golang.org/x/exp/inotify"
)

type process func(string, string)

func main() {
	// Watch all of the arguments passed in

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	messages := make(chan bool)
	control := make(chan bool, 1)
	numWatches := len(os.Args) - 1

	for _, arg := range os.Args[1:] {
		log.Println("Arg:", arg)
		go func() {
			watchItem(arg, squareVideo, control)
			messages <- true
		}()
	}

	log.Println("Starting watches")
	for messages != nil && sigs != nil {
		select {
		case <-messages:
			numWatches--
			if numWatches <= 0 {
				close(messages)
				messages = nil
				log.Println("All watches have returned")
			}
		case <-sigs:
			control <- true
			close(sigs)
			sigs = nil
			log.Println("Received event, cleaning up")
		}
	}
	close(control)
	log.Println("Exiting program")
}

func watchItem(path string, fn process, control chan bool) {
	log.Println("Watching:" + path)
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
			processEvent(ev, fn)
		case err := <-watcher.Error:
			log.Println("error:", err)
		case <-control:
			control <- true
			return
		}
	}
}

func processEvent(ev *inotify.Event, fn process) {
	if ev.Mask&inotify.IN_CLOSE_WRITE != 0 {
		log.Println("The file " + ev.Name + " has been closed")
		matched, err := regexp.MatchString("(?i).*rectangle\\.mp4$", ev.Name)
		if err != nil {
			log.Fatal("Failed to match regex", err)
		}
		if matched {
			re := regexp.MustCompile("(?i)rectangle\\.mp4$")
			output := re.ReplaceAllString(ev.Name, "Square.mp4")
			fn(ev.Name, output)
		}
	}
}

func squareVideo(input string, output string) {
	log.Println("Processing video", input, output)
	cmd := exec.Command("ffmpeg", "-y", "-i", input, "-filter:v", "crop=in_h:in_h", "-c:a", "copy", output)
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Printlnl("Unable to process video: ", string(result), err)
		return
	}
	log.Println("Successfully cropped video to " + output)
}
