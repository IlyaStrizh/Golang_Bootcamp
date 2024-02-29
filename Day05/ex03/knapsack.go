package ex03

import (
	heap "Day05/ex02"
)

func GrabPresents(presents []heap.Present, m int) []heap.Present {
	w := len(presents)
	dp := make([][]int, w+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	for i := 1; i <= w; i++ {
		for j := 1; j <= m; j++ {
			if presents[i-1].Size <= j {
				dp[i][j] = max(presents[i-1].Value+dp[i-1][j-presents[i-1].Size], dp[i-1][j])
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}

	result := make([]heap.Present, 0)
	for i, j := w, m; i > 0 && j > 0; i-- {
		if dp[i][j] != dp[i-1][j] {
			// Делает "result = append(result, presents[i-1])" - наоборот.
			// Создает новый слайс, содержащий элемент "presents[i-1]"" в начале,
			// а затем все элементы из "result" распаковываются и добавляются после этого элемента в новом слайсе
			result = append([]heap.Present{presents[i-1]}, result...)
			j -= presents[i-1].Size
		}
	}

	return result
}
