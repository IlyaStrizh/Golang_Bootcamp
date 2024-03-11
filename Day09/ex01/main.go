package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	urlChannel := make(chan string)

	go func() {
		for i := 0; i < 250; i++ {
			urlChannel <- "https://www.example.com"
		}
		close(urlChannel)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		fmt.Println("Crawling stopping...")
		cancel()
	}()

	html := crawlWeb(urlChannel, ctx)
	for h := range html {
		fmt.Println(*h)
	}
}
