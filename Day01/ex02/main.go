package main

import (
	"day01/ex02/readerTXT"
	"flag"
	"fmt"
	"os"
)

func main() {
	pathOld := flag.String("old", "", "Первый дамп (txt)")
	pathNew := flag.String("new", "", "Второй дамп (txt)")
	flag.Parse()

	if *pathOld == "" || *pathNew == "" {
		fmt.Println("Необходимо указать два файла")
		os.Exit(0)
	}

	old := readerTXT.NewTXT()
	old.Read(pathOld)
	new := readerTXT.NewTXT()
	new.Read(pathNew)
	old.Compare(new)
}
