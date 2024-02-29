package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestSendRequests(t *testing.T) {
	t.Run("LimitTest", func(t *testing.T) {
		status429Counter := 0
		status200Counter := 0

		for {
			for j := 0; j < 200; j++ {
				response, err := http.Get("http://localhost:8888")
				if err != nil {
					fmt.Println("Error making request:", err)
					continue
				}

				if response.StatusCode == http.StatusTooManyRequests {
					status429Counter++
				} else if response.StatusCode == http.StatusOK {
					status200Counter++
				}
			}

			fmt.Printf("Got %d http.StatusTooManyRequests [429] | Got %d http.StatusOK [200]\n", status429Counter, status200Counter)
			status429Counter = 0
			status200Counter = 0

			time.Sleep(time.Second)
		}
	})
}
