package main

import (
	"fmt"
	"os"
)

func generate(path string) {
	currentPkg, err := NewTree(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(currentPkg.GenerateFile())
}
