package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	args := os.Args[1:]

	for scanner.Scan() {
		line := scanner.Text()
		commandArgs := append(args, strings.Fields(line)...)

		cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Ошибка выполнения команды:", err)
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Ошибка чтения ввода:", scanner.Err())
	}
}
