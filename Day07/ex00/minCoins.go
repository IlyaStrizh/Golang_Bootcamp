package ex00

import "sort"

func MinCoins(val int, coins []int) []int {
	res := make([]int, 0)
	i := len(coins) - 1
	for i >= 0 {
		for val >= coins[i] {
			val -= coins[i]
			res = append(res, coins[i])
		}
		i -= 1
	}
	return res
}

func MinCoins2(val int, coins []int) []int {
	if val <= 0 || len(coins) == 0 || coins == nil {
		return []int{}
	}
	// Создаем массив для хранения минимального количества монет в каждой сумме
	dp := make([]int, val+1)
	// Заполняем массив dp с использованием динамического программирования
	for i := 1; i <= val; i++ {
		dp[i] = val + 1
		for _, coin := range coins {
			if coin <= i && coin > 0 && dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
			}
		}
	}
	// Возвращаем пустой срез, если невозможно получить сумму val с использованием заданных монет
	if dp[val] == val+1 {
		return []int{}
	}

	res := make([]int, 0)
	for i := val; i > 0; {
		for _, coin := range coins {
			if coin <= i && coin > 0 && dp[i-coin]+1 == dp[i] {
				res = append(res, coin)
				i -= coin
				break
			}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] > res[j]
	})

	return res
}
