package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a command line argument!")
		os.Exit(1)
	}
	for _, path := range filepath.SplitList(os.Getenv("PATH")) {
		fullPath := filepath.Join(path, os.Args[1])
		_, err := os.Stat(fullPath)
		if !os.IsNotExist(err) {
			fmt.Println(os.Args[1], "->", fullPath)
			return
		}
	}
	os.Exit(1)
}
