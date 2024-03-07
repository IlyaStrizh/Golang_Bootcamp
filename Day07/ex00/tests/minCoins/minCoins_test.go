package tests

import (
	"Day07/ex00"
	"reflect"
	"testing"

	"github.com/fatih/color"
)

func TestMinCoins(t *testing.T) {
	var (
		val      int
		coins    []int
		expected []int
		result   []int
	)

	t.Run("Test#1", func(t *testing.T) {
		val = 13
		coins = []int{1, 5, 10, 15, 50}
		expected = []int{10, 1, 1, 1}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Errorf(errorMsg)
		}
	})

	t.Run("Test#2", func(t *testing.T) {
		val = 13
		coins = []int{5, 10, 1}
		expected = []int{10, 1, 1, 1}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#3", func(t *testing.T) {
		val = 10
		coins = nil
		expected = []int{}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#4", func(t *testing.T) {
		val = 15
		coins = []int{1, 10, 5}
		expected = []int{10, 5}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#5", func(t *testing.T) {
		val = 15
		coins = []int{1, 10, 5, 1, 10, 5}
		expected = []int{10, 5}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#6", func(t *testing.T) {
		val = 14
		coins = []int{1, 2, 7, 10}
		expected = []int{7, 7}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#7", func(t *testing.T) {
		val = -14
		coins = []int{1, 2, 7, 10}
		expected = []int{}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#8", func(t *testing.T) {
		val = 0
		coins = []int{1, 2, 7, 10}
		expected = []int{}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#9", func(t *testing.T) {
		val = 6
		coins = []int{2, 4, 5, 7, 10}
		expected = []int{4, 2}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})

	t.Run("Test#10", func(t *testing.T) {
		val = -6
		coins = []int{-1, -4, -5, 4, 2, 7, 10}
		expected = []int{}
		result = []int{}

		result = ex00.MinCoins(val, coins)
		if !reflect.DeepEqual(result, expected) {
			errorMsg := color.MagentaString(`minCoins: "%v", want: "%v"`, result, expected)
			t.Error(errorMsg)
		}
	})
}
