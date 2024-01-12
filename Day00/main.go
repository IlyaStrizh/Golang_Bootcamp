package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type Flag struct {
	mean, median, mode, sd bool
}

func initFlag() Flag {
	var f Flag

	flag.BoolVar(&f.mean, "mean", false, "calculate mean")
	flag.BoolVar(&f.median, "median", false, "calculate median")
	flag.BoolVar(&f.mode, "mode", false, "calculate mode")
	flag.BoolVar(&f.sd, "sd", false, "calculate standard deviation")
	flag.Parse()

	return f
}

func parseInput() []float64 {
	numbers := []float64{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		num, err := strconv.ParseFloat(scanner.Text(), 64)
		if num < -100000 || num > 100000 {
			fmt.Fprintln(os.Stderr, "Число > 100000 или < -100000")
			continue
		}
		if err == nil {
			numbers = append(numbers, num)
		} else {
			fmt.Fprintln(os.Stderr, "Ошибка ввода", err)
			continue
		}
	}
	if len(numbers) == 0 {
		fmt.Fprintln(os.Stderr, "Пустой ввод")
		os.Exit(1)
	}
	sort.Float64s(numbers)

	return numbers
}

func printResults(nums []float64, f Flag) {
	if !f.mean && !f.median && !f.mode && !f.sd {
		f.mean = true
		f.median = true
		f.mode = true
		f.sd = true
	}

	if f.mean {
		fmt.Printf("Mean: %.2f\n", calcMean(nums))
	}
	if f.median {
		fmt.Printf("Median: %.2f\n", calcMedian(nums))
	}
	if f.mode {
		fmt.Printf("Mode: %.2f\n", calcMode(nums))
	}
	if f.sd {
		fmt.Printf("SD: %.2f\n", calcSD(nums, calcMean(nums)))
	}
}

func calcMean(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

func calcMedian(numbers []float64) float64 {
	n := len(numbers)
	if n%2 == 1 {
		return float64(numbers[n/2])
	}
	return float64(numbers[n/2-1]+numbers[n/2]) / 2.0
}

func calcMode(numbers []float64) float64 {
	counts := make(map[float64]int)
	for _, num := range numbers {
		counts[num]++
	}
	maxCount := 0
	mode := math.Inf(1)
	for num, count := range counts {
		if count > maxCount { //|| (count == maxCount && num < mode) {
			maxCount = count
			mode = num
		}
	}
	if maxCount == 1 {
		mode = 0.00
	}

	return mode
}

func calcSD(numbers []float64, mean float64) float64 {
	sumSquares := 0.0
	for _, num := range numbers {
		diff := num - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(numbers))
	return math.Sqrt(variance)
}

func main() {
	f := initFlag()
	printResults(parseInput(), f)
}
