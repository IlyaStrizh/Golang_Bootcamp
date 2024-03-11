package ex00

import (
	"sync"
	"time"
)

func sleepSort(arr []int) <-chan int {
	len := len(arr)

	ch := make(chan int, len)
	wg := new(sync.WaitGroup)

	wg.Add(len)
	for i := range arr {
		n := arr[i]
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(n) * time.Second)
			ch <- n
		}()
	}
	wg.Wait()
	close(ch)

	return ch
}
