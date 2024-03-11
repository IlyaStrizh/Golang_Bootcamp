package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func crawlWeb(urlChannel <-chan string, ctx context.Context) <-chan *string {
	resultChannel := make(chan *string)

	go runCrawling(urlChannel, ctx, resultChannel)

	return resultChannel
}

func runCrawling(urlChannel <-chan string, ctx context.Context, resultChannel chan<- *string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 8) // Ограничение количества горутин
	defer close(resultChannel)

	for url := range urlChannel {
		select {
		case <-ctx.Done():
			return
		case semaphore <- struct{}{}: // Берем "токен" из семафора
			wg.Add(1)
			go crawl(url, &wg, semaphore, resultChannel)
		}
	}

	wg.Wait()
}

func crawl(url string, wg *sync.WaitGroup, s chan struct{}, res chan<- *string) {
	defer wg.Done()
	defer func() { <-s }() // Возвращаем "токен" в семафор

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	result := string(body)
	res <- &result
}
