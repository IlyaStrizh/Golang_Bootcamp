package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
)

type flags struct {
	l, m, w bool
}

func initFlags() *flags {
	f := &flags{}

	flag.BoolVar(&f.l, "l", false, "Подсчет строк")
	flag.BoolVar(&f.m, "m", false, "Подсчет букв")
	flag.BoolVar(&f.w, "w", false, "Подсчет слов")
	flag.Parse()

	return f
}

func lineCount(scanner *bufio.Scanner) uint64 {
	var count uint64
	for scanner.Scan() {
		count++
	}
	return count
}

func wordCount(scanner *bufio.Scanner) uint64 {
	var count uint64
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		count++
	}
	return count
}

func charCount(scanner *bufio.Scanner) uint64 {
	var count uint64
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		count++
	}
	return count
}

func printResult(path string, f *flags, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка открытия файла", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	switch {
	case f.l:
		fmt.Printf("%d\t%s\n", lineCount(scanner), path)
	case f.m:
		fmt.Printf("%d\t%s\n", charCount(scanner), path)
	case f.w:
		fmt.Printf("%d\t%s\n", wordCount(scanner), path)
	}
}

func main() {
	fl := initFlags()
	count := flag.NFlag()

	if count == 0 {
		fl.w = true
	} else if count > 1 {
		fmt.Fprintln(os.Stderr, "Укажите не более одного флага")
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)
	if len(flag.Args()) > 0 {
		paths := flag.Args()
		for _, path := range paths {
			wg.Add(1)
			go printResult(path, fl, wg)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Укажите минимум один путь к файлу")
		os.Exit(1)
	}
	wg.Wait()
}
