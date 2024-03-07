package main

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
// #include "application.h"
// #include "window.h"
import "C"
import "unsafe"

func main() {
	// Инициализация приложения
	C.InitApplication()

	title := C.CString("School 21")
	defer C.free(unsafe.Pointer(title))

	// Создание окна
	wndPtr := C.Window_Create(500, 500, 300, 200, title)

	// Сделать окно активным и отображаемым
	C.Window_MakeKeyAndOrderFront(wndPtr)

	// Запуск приложения
	C.RunApplication()
}
