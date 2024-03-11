package ex00

import (
	"reflect"
	"testing"

	"github.com/fatih/color"
)

func TestSleepSort(t *testing.T) {
	var (
		arr      []int
		expected []int
	)

	t.Run("Test#1", func(t *testing.T) {
		arr = []int{5, 3, 9, 8, 1, 7}
		expected = []int{1, 3, 5, 7, 8, 9}
		var result []int

		resCh := sleepSort(arr)
		for i := range resCh {
			result = append(result, i)
		}
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.HiMagentaString(`sleepSort: "%v", want: "%v"`, result, expected)
			t.Errorf(errorMsg)
		}
	})

	t.Run("Test#2", func(t *testing.T) {
		arr = []int{}
		var expected []int
		var result []int

		resCh := sleepSort(arr)
		for i := range resCh {
			result = append(result, i)
		}
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.HiMagentaString(`sleepSort: "%v", want: "%v"`, result, expected)
			t.Errorf(errorMsg)
		}
	})

	t.Run("Test#3", func(t *testing.T) {
		arr = nil
		var expected []int
		var result []int

		resCh := sleepSort(arr)
		for i := range resCh {
			result = append(result, i)
		}
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.HiMagentaString(`sleepSort: "%v", want: "%v"`, result, expected)
			t.Errorf(errorMsg)
		}
	})
}
