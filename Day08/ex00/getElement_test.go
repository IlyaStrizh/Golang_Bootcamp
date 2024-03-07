package ex00

import (
	"reflect"
	"testing"

	"github.com/fatih/color"
)

func TestGetElement(t *testing.T) {
	var (
		idx      int
		arr      []int
		expected int
	)

	t.Run("Test#1", func(t *testing.T) {
		arr = []int{1, 5, 10, 15, 50}

		for i := range arr {
			result, err := getElement(arr, i)
			expected = arr[i]
			if err != nil {
				errorMsg := color.MagentaString(`Error: "%s"`, err)
				t.Errorf(errorMsg)
			}
			if !reflect.DeepEqual(result, expected) {
				errorMsg := color.MagentaString(`getElement: "%v", want: "%v"`, result, expected)
				t.Errorf(errorMsg)
			}
		}
	})

	t.Run("Test#2", func(t *testing.T) {
		idx = -1
		arr = []int{1, 5, 10, 15, 50}
		expected = 5

		result, err := getElement(arr, idx)
		if err == nil {
			errorMsg := color.MagentaString(`getElement: "%v", want: Error`, result)
			t.Errorf(errorMsg)
		}
	})

	t.Run("Test#3", func(t *testing.T) {
		idx = 5
		arr = []int{1, 5, 10, 15, 50}
		expected = 5

		result, err := getElement(arr, idx)
		if err == nil {
			errorMsg := color.MagentaString(`getElement: "%v", want: Error`, result)
			t.Errorf(errorMsg)
		}
	})

	t.Run("Test#4", func(t *testing.T) {
		idx = 3
		arr = nil
		expected = 5

		result, err := getElement(arr, idx)
		if err == nil {
			errorMsg := color.MagentaString(`getElement: "%v", want: Error`, result)
			t.Errorf(errorMsg)
		}
	})

	t.Run("Test#5", func(t *testing.T) {
		idx = 3
		arr = []int{}
		expected = 5

		result, err := getElement(arr, idx)
		if err == nil {
			errorMsg := color.MagentaString(`getElement: "%v", want: Error`, result)
			t.Errorf(errorMsg)
		}
	})
}
