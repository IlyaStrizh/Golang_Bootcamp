package main

import (
	"day01/ex00/reader"
	"flag"
	"fmt"
	"os"
)

func main() {
	pathOld := flag.String("old", "", "Первый рецепт (json или xml)")
	pathNew := flag.String("new", "", "Второй рецепт (json или xml)")
	flag.Parse()

	if *pathOld == "" || *pathNew == "" {
		fmt.Println("Необходимо указать два файла")
		os.Exit(0)
	}

	old, err := reader.NewDBReader(pathOld)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	new, err := reader.NewDBReader(pathNew)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	old.Read(pathOld)
	new.Read(pathNew)
	old.Compare(new)
}
