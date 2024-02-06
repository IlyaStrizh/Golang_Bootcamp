package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type flags struct {
	sl, d, f bool
	ext      string
	path     string
}

func initFlags() *flags {
	f := &flags{}

	flag.BoolVar(&f.sl, "sl", false, "Поиск символических ссылок")
	flag.BoolVar(&f.d, "d", false, "Поиск каталогов")
	flag.BoolVar(&f.f, "f", false, "Поиск файлов")
	flag.StringVar(&f.ext, "ext", "", `Поиск по расширению (только с флагом "-f")`)

	flag.Parse()

	if len(flag.Args()) > 0 {
		f.path = flag.Arg(0)
	} else {
		fmt.Fprintln(os.Stderr, "Не указан путь для поиска файлов")
		os.Exit(1)
	}

	return f
}

func findFiles(f *flags) {
	if !f.f && !f.d && !f.sl {
		f.f, f.d, f.sl = true, true, true
	}

	err := filepath.Walk(f.path, func(path string, fileInfo os.FileInfo, err error) error {
		if os.IsPermission(err) {
			return filepath.SkipDir
		} else if err != nil {
			return err
		}

		if f.ext != "" {
			fileExt := filepath.Ext(path)
			if fileInfo.Mode().IsRegular() && fileExt == "."+f.ext {
				fmt.Println(path)
			}
		} else if f.f && fileInfo.Mode().IsRegular() {
			fmt.Println(path)
		}
		if f.sl {
			if symlink, _ := filepath.EvalSymlinks(path); symlink != path {
				if _, errExist := os.Stat(symlink); errExist == nil {
					fmt.Println(path, "->", symlink)
				} else {
					fmt.Println(path, "->", "[broken]")
				}
			}
		}
		if f.d && fileInfo.IsDir() {
			fmt.Println(path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка поиска файлов", err)
		os.Exit(1)
	}
}

func main() {
	f := initFlags()

	if f.ext != "" && !f.f {
		fmt.Fprintln(os.Stderr, `Флаг "-ext" используется с флагом "-f"`)
		os.Exit(1)
	}

	findFiles(f)
}
