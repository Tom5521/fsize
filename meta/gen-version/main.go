package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var buf bytes.Buffer

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	str := buf.String()
	str = strings.ReplaceAll(str, "\x0a", "")

	err = os.WriteFile("version.txt", []byte(str), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}