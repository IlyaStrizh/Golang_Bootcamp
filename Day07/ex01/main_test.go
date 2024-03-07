// go test -bench . -benchmem
// go test -bench=. -benchmem -cpuprofile cpu.pprof
//go tool pprof исполняемый_файл cpu.pprof

package ex01

import (
	"Day07/ex00"
	"testing"

	"github.com/pkg/profile"
)

func BenchmarkMinCoinsPPROF(b *testing.B) {
	defer profile.Start(profile.ProfilePath(".")).Stop()

	val := 100
	b.ResetTimer()
	b.ReportAllocs()
	coins := []int{1, 5, 10, 25}
	for i := 0; i < b.N; i++ {
		ex00.MinCoins2(val, coins)
	}
}

func BenchmarkMinCoins(b *testing.B) {
	val := 10000
	b.ResetTimer()
	b.ReportAllocs()
	coins := []int{1, 3, 5, 7}
	for i := 0; i < b.N; i++ {
		ex00.MinCoins(val, coins)
	}
}

func BenchmarkMinCoins2(b *testing.B) {
	val := 10000
	b.ResetTimer()
	b.ReportAllocs()
	coins := []int{1, 3, 5, 7}
	for i := 0; i < b.N; i++ {
		ex00.MinCoins2(val, coins)
	}
}
