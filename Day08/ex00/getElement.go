package ex00

import (
	"errors"
	"unsafe"
)

func getElement(arr []int, idx int) (int, error) {
	l := len(arr)
	if idx < 0 || idx >= l || arr == nil {
		return 0, errors.New("некорректные входные данные")
	}

	slice := unsafe.Pointer(&arr[0])
	elementPtr := (*int)(unsafe.Pointer(uintptr(slice) + uintptr(idx)*unsafe.Sizeof(arr[0])))

	return *elementPtr, nil
}
