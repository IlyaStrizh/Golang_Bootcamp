package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func initFlag() *string {
	f := flag.String("a", "", "Путь для сохранения архивов")
	flag.Parse()

	return f
}

func checkFolder(archivePath *string) {
	info, err := os.Stat(*archivePath)
	if !os.IsPermission(err) {
		if os.IsNotExist(err) || !info.Mode().IsDir() {
			err := os.Mkdir(*archivePath, 0755)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Не удалось создать директорию:", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Отказано в доступе к %s\n:", *archivePath)
		os.Exit(1)
	}
}

func checkFile(filePath *string) {
	info, err := os.Stat(*filePath)
	if !os.IsPermission(err) {
		if os.IsNotExist(err) || !info.Mode().IsRegular() {
			fmt.Fprintf(os.Stderr, "Файл %s не существует или не является файлом\n", *filePath)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Отказано в доступе к %s:\n", *filePath)
		os.Exit(1)
	}
}

func archiveFile(filePath string, archivePath string) error {
	if archivePath == "" {
		archivePath = filepath.Dir(filePath)
	}
	baseName := filepath.Base(filePath)
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	newName := fmt.Sprintf("%s_%d.tar.gz", name, time.Now().Unix())

	tarFile, err := os.Create(filepath.Join(archivePath, newName))
	if err != nil {
		return err
	}
	defer tarFile.Close()
	gw := gzip.NewWriter(tarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	fileToArchive, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToArchive.Close()
	fileInfo, err := fileToArchive.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	if err != nil {
		return err
	}
	header.Name = baseName
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(tw, fileToArchive); err != nil {
		return err
	}

	return nil
}

func main() {
	archivePath := initFlag()
	if flag.NFlag() == 1 {
		checkFolder(archivePath)
	}

	wg := new(sync.WaitGroup)
	for _, filePath := range flag.Args() {
		checkFile(&filePath)
		errChan := make(chan error, 1)
		wg.Add(1)

		go func() {
			defer wg.Done()
			errChan <- archiveFile(filePath, *archivePath)
		}()

		wg.Wait()
		if err := <-errChan; err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка архивирования:", err)
			os.Exit(1)
		}
	}
}
