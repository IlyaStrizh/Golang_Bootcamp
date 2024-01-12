package main

import (
	"math"
	"testing"
)

func TestCalcMean(t *testing.T) {
	numbers := []float64{1, 2, 3, 4, 5}
	want := 3.0
	if got := calcMean(numbers); got != want {
		t.Errorf("calcMean() = %.2f, want %.2f", got, want)
	}
}

func TestCalcMedian(t *testing.T) {
	numbers := []float64{1, 2, 3, 4, 5}
	want := 3.0
	if got := calcMedian(numbers); got != want {
		t.Errorf("calcMedian() = %.2f, want %.2f", got, want)
	}
}

func TestCalcMode(t *testing.T) {
	numbers := []float64{1, 2, 2, 3, 4}
	want := 2.0
	if got := calcMode(numbers); got != want {
		t.Errorf("calcMode() = %.2f, want %.2f", got, want)
	}
}

func TestCalcSD(t *testing.T) {
	numbers := []float64{1, 2, 3, 4, 5}
	want := 1.41421356
	if got := math.Round(calcSD(numbers, calcMean(numbers))*1e8) / 1e8; got != want {
		t.Errorf("calcSD() = %.8f, want %.8f", got, want)
	}
}
