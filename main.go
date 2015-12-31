package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/tillberg/ansi-log"
)

var colors []string = []string{"red", "yellow", "magenta", "cyan", "blue"}

func writeStuff(num int, stop chan bool, done chan bool) {
	i := 0
	color := colors[num/4]
	out := alog.New(os.Stderr, fmt.Sprintf("@(dim)[@(green:writer-%d)] ", num), alog.Lelapsed)
	format := fmt.Sprintf("@(dim:My number is) @(%s:%%d)", color)
	for {
		select {
		case <-stop:
			out.Close()
			done <- true
			return
		default:
			out.Replacef(format, i)
			i += num
			if i > 10000 {
				i = 0
				out.Println()
			}
		}
		time.Sleep(time.Microsecond)
	}
}

func main() {
	alog.EnableMultilineMode()
	alog.EnableColorTemplate()
	stop := make(chan bool)
	done := make(chan bool)

	numLines := 20
	for n := 0; n < numLines; n++ {
		go writeStuff(n, stop, done)
	}
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			alog.Println("Received SIGINT. Shutting down...")
			for n := 0; n < numLines; n++ {
				stop <- true
			}
			for n := 0; n < numLines; n++ {
				<-done
			}
			alog.Println("Exiting...")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
