/*
Генерируем документацию:
	1. Установка godoc: $ go install golang.org/x/tools/cmd/godoc
	2. Генерация: $ ~/go/bin/godoc -http=:6060
	3. Переходим по адресу http://localhost:6060/pkg/day07/ex02/ для просмотра документации
	4. Я остановился на варианте сохранения HTML-страницы в файл:
		curl -s -o docs.html http://localhost:6060/pkg/Day07/ex02/
	5. Архивируем: zip docs.zip docs.html
*/

package ex02

import "sort"

/*
MinCoins - функция, которая приведена в качестве примера. В некоторых случаях работает некорректно.

Пример:
val = 15, coins = []int{1, 10, 5}
minCoins_test.go:67: minCoins: "[5 5 5]", want: "[10 5]".

MinCoins использует жадный алгоритм (берет монету наибольшего номинала, если она не превышает значение val).

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
*/
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

/*
MinCoins2 - новая версия MinCoins.

Она обрабатывает случаи, с которыми не справляется первая версия.

MinCoins2 реализован при помощи динамического программирования.
Он создает массив сумм, содержащий минимальное количество монет для каждой суммы от 1 до val.

Алгоритм в функции MinCoins2 гарантирует, что найдено оптимальное решение.
Если оптимального решения не существует - вернет пустой фрагмент.

	func MinCoins2(val int, coins []int) []int {
		// Проверяем правильность ввода
		if val <= 0 || len(coins) == 0 || coins == nil {
			return []int{}
		}
		// Создаем массив для хранения минимального количества монет в каждой сумме
		dp := make([]int, val+1)
		// Заполняем массив dp с использованием динамического программирования,
		// создавая массив сумм для каждой суммы от 1 до val
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
				// Этот цикл восстанавливает оптимальный набор монет для заданной суммы val.
				// Мы начинаем с val и на каждой итерации выбираем монету "coin", которая привела
				// к уменьшению минимального количества монет для текущей суммы val
				if coin <= i && coin > 0 && dp[i-coin]+1 == dp[i] {
					res = append(res, coin)
					i -= coin
					break
				}
			}
		}
		// Сортируем результат
		sort.Slice(res, func(i, j int) bool {
			return res[i] > res[j]
		})

		return res
	}
*/
func MinCoins2(val int, coins []int) []int {
	// Проверяем правильность ввода
	if val <= 0 || len(coins) == 0 || coins == nil {
		return []int{}
	}
	// Создаем массив для хранения минимального количества монет в каждой сумме
	dp := make([]int, val+1)
	// Заполняем массив dp с использованием динамического программирования,
	// создавая массив сумм для каждой суммы от 1 до val
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
			// Этот цикл восстанавливает оптимальный набор монет для заданной суммы val.
			// Мы начинаем с val и на каждой итерации выбираем монету "coin", которая привела
			// к уменьшению минимального количества монет для текущей суммы val
			if coin <= i && coin > 0 && dp[i-coin]+1 == dp[i] {
				res = append(res, coin)
				i -= coin
				break
			}
		}
	}
	// Сортируем результат
	sort.Slice(res, func(i, j int) bool {
		return res[i] > res[j]
	})

	return res
}
