package main

import (
	"fmt"
	"sort"
	"strings"
)

type Tree struct {
	root  string
	nodes map[string]*Node
}

type Node struct {
	capabilities *Set
	imports      *Set
}

func NewTree(start string) (*Tree, error) {
	t := &Tree{nodes: map[string]*Node{}}

	rootPackagePath, err := parsePackagePath(start)
	if err != nil {
		return nil, err
	}

	if err := t.addPackagesRecursive(rootPackagePath); err != nil {
		return nil, err
	}

	t.root = rootPackagePath
	return t, nil
}

func (t *Tree) addPackagesRecursive(pkgPath string) (err error) {
	if _, ok := t.nodes[pkgPath]; ok {
		return nil
	}

	if pkgPath == "C" {
		t.addNote("C", nil, []string{"syscall"})
		return nil
	}

	imports, err := parsePackage(pkgPath)
	if err != nil {
		return err
	}

	var external []string
	var standard []string
	for _, im := range imports {
		if !isStandard(im) {
			if err := t.addPackagesRecursive(im); err != nil {
				return err
			}

			external = append(external, im)
		} else {
			standard = append(standard, im)
		}
	}

	t.addNote(pkgPath, external, toCapabilities(standard))
	return nil
}

func isStandard(name string) bool {
	_, ok := standardPackages[name]
	return ok
}

func (t *Tree) addNote(name string, imports, capabilities []string) {
	t.nodes[name] = &Node{
		capabilities: NewSet(capabilities...),
		imports:      NewSet(imports...),
	}
}

func (t *Tree) GenerateFile() string {
	sb := strings.Builder{}
	sb.WriteString(goCapLine(t.root, t.nodes[t.root].capabilities.Values()))
	sb.WriteString("\n")

	for _, name := range sortedKeys(t.nodes) {
		if name != t.root {
			values := t.nodes[name].capabilities.Values()
			if len(values) > 0 {
				sb.WriteString(goCapLine(name, t.nodes[name].capabilities.Values()))
			}
		}
	}

	return sb.String()
}

func sortedKeys(m map[string]*Node) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func goCapLine(name string, capabilities []string) string {
	return fmt.Sprintf("%s (%s)\n", name, strings.Join(capabilities, ", "))
}
