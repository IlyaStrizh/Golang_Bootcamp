package ex02

import (
	"container/heap"
	"errors"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (p PresentHeap) Len() int {
	return len(p)
}

func (p PresentHeap) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		return p[i].Size < p[j].Size
	}
	return p[i].Value > p[j].Value
}

func (p *PresentHeap) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *PresentHeap) Push(x interface{}) {
	*p = append(*p, x.(Present))
}

func (p *PresentHeap) Pop() interface{} {
	n := len(*p)
	x := (*p)[n-1]  // Получаем последний элемент
	*p = (*p)[:n-1] // Удаляем последний элемент

	return x
}

func GetNCoolestPresents(p []Present, n int) ([]Present, error) {
	if n < 0 || n > len(p) {
		return []Present{}, errors.New(`Invalid value "n"`)
	}

	tmp := make(PresentHeap, len(p))
	for i, present := range p {
		tmp[i].Value = present.Value
		tmp[i].Size = present.Size
	}
	heap.Init(&tmp)

	result := make([]Present, n)
	for i := range result {
		result[i] = heap.Pop(&tmp).(Present)
	}

	return result, nil
}
