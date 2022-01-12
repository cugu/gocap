package main

import "sort"

var permissionMap = map[string][]string{
	"os":        {"file"},
	"io/ioutil": {"file"},
	"net":       {"network"},
	"net/http":  {"network"},
	"os/exec":   {"execute"},
	"syscall":   {"syscall"},
	"reflect":   {"reflect"},
	"unsafe":    {"unsafe"},
}

func toCapabilities(dependencies []string) []string {
	capabilities := NewSet()
	for _, dependency := range dependencies {
		if capas, ok := permissionMap[dependency]; ok {
			capabilities.Add(capas...)
		}
	}
	cs := capabilities.Values()
	sort.Strings(cs)
	return cs
}
