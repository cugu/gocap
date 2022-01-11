package main

import (
	"os/exec"
	"strings"
)

var standardPackages = make(map[string]struct{})

func init() {
	p, err := exec.Command("go", "list", "std").Output()
	if err != nil {
		panic(err)
	}

	for _, p := range strings.Fields(string(p)) {
		standardPackages[p] = struct{}{}
	}
}
