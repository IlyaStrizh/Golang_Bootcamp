package ex02

import "sync"

func multiplex(ch ...<-chan interface{}) <-chan interface{} {
	n := len(ch)
	res := make(chan interface{}, n)
	wg := new(sync.WaitGroup)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			for j := range ch[i] {
				res <- j
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
