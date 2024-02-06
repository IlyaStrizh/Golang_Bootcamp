package readerTXT

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TXT map[string]struct{}

func NewTXT() *TXT {
	return &TXT{}
}

type TXTreader interface {
	Read(s *string)
	Compare(old *TXT)
}

func (t *TXT) Read(s *string) {
	fileExt := filepath.Ext(*s)
	if fileExt == ".txt" {
		file, err := os.Open(*s)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка открытия файла:", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			(*t)[line] = struct{}{}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка чтения в буфер:", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Неподдерживаемый тип файла")
		os.Exit(1)
	}
}

func (t *TXT) Compare(new *TXT) {
	for key := range *new {
		if _, ok := (*t)[key]; !ok && strings.TrimSpace(key) != "" {
			fmt.Printf("ADDED %s\n", key)
		}
	}
	for key := range *t {
		if _, ok := (*new)[key]; !ok && strings.TrimSpace(key) != "" {
			fmt.Printf("REMOVED %s\n", key)
		}
	}
}
