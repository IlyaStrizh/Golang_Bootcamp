package main

import (
	"day01/ex00/reader"
	"flag"
	"fmt"
	"os"
)

func main() {
	path := flag.String("f", "", "Файл json или xml")
	flag.Parse()
	db, err := reader.NewDBReader(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	db.Read(path)
	fmt.Println(db.Convert())
}
