package ex02

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/fatih/color"
)

func TestSleepSort(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})
	ch4 := make(chan interface{})
	ch5 := make(chan interface{})

	go func() {
		for i := 1; i <= 10; i++ {
			ch1 <- i
			ch2 <- i + 10
			ch3 <- i + 20
			ch4 <- i + 30
			ch5 <- i + 40
		}
		close(ch1)
		close(ch2)
		close(ch3)
		close(ch4)
		close(ch5)
	}()

	resCh := multiplex(ch2, ch1, ch3, ch5, ch4)

	count := 0
	var result []interface{}
	for i := range resCh {
		count++
		result = append(result, i)
	}

	sort.Slice(result, func(i, j int) bool { return result[i].(int) < result[j].(int) })
	expected := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50}

	if count == 50 && reflect.DeepEqual(result, expected) {
		testMsg := color.HiGreenString(">>> %d goroutines were proceeded\n", count)
		fmt.Fprintln(os.Stdout, testMsg)
	} else {
		errorMsg := color.HiMagentaString(">>> %d goroutines were proceeded\n multiplex:\n \"%v\"\n want:\n \"%v\"", count, result, expected)
		t.Errorf(errorMsg)
	}
}
