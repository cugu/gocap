package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
)

type UnnecessaryPackage struct {
	Package string
}

func (e *UnnecessaryPackage) Error() string {
	return fmt.Sprintf("unnecessary package '%s', please remove from go.cap file", e.Package)
}

type PackageNotProvided struct {
	Package string
}

func (e *PackageNotProvided) Error() string {
	return fmt.Sprintf("package '%s' not listed in go.cap file, please add to go.cap file", e.Package)
}

type UnnecessaryCapability struct {
	Capability string
}

func (e *UnnecessaryCapability) Error() string {
	return fmt.Sprintf("unnecessary capability '%s', please remove from go.cap file", e.Capability)
}

type CapabilityNotProvided struct {
	Capability string
}

func (e *CapabilityNotProvided) Error() string {
	return fmt.Sprintf("capability '%s' not provided by go.cap file, add to go.cap file if you want to grant the capability", e.Capability)
}

func checkCmd(path string, ignore string) {
	file, err := parseGoCap(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tree, err := NewTree(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if ignore != "" {
		filterTree(tree, ignore)
		file.External = filterFile(file, ignore)
	}

	packageErrors := check(tree, file)

	if len(packageErrors) > 0 {
		if rootErrs, ok := packageErrors[tree.root]; ok {
			fmt.Println(tree.root)
			for _, e := range rootErrs {
				fmt.Println("\t" + e.Error())
			}
		}
		for _, name := range sortedPackageErrorKeys(packageErrors) {
			if name == tree.root {
				continue
			}
			fmt.Println(name)
			for _, e := range packageErrors[name] {
				fmt.Println("\t" + e.Error())
			}
		}
		os.Exit(1)
	}
}

func filterTree(tree *Tree, ignore string) {
	for name := range tree.nodes {
		match, err := regexp.MatchString(ignore, name)
		if err != nil {
			continue
		}
		if match {
			delete(tree.nodes, name)
		}
	}
}

func filterFile(file *File, ignore string) []Line {
	var external []Line
	for _, line := range file.External {
		match, err := regexp.MatchString(ignore, line.Package)
		if err != nil {
			continue
		}
		if !match {
			external = append(external, line)
		}
	}
	return external
}

func sortedPackageErrorKeys(m map[string][]error) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func check(tree *Tree, file *File) map[string][]error {
	packageErrors := map[string][]error{}
	for name := range tree.nodes {
		errors := checkPackage(name, tree, file)

		if len(errors) > 0 {
			packageErrors[name] = errors
		}
	}

	for _, line := range file.External {
		if _, ok := tree.nodes[line.Package]; !ok {
			packageErrors[line.Package] = []error{&UnnecessaryPackage{Package: line.Package}}
		}
	}

	return packageErrors
}

func checkPackage(name string, tree *Tree, file *File) []error {
	values := tree.nodes[name].capabilities.Values()
	if _, ok := file.expected[name]; ok {
		return checkGoCapLine(file.expected[name], values)
	}

	if len(values) > 0 {
		return []error{&PackageNotProvided{Package: name}}
	}
	return nil
}

func checkGoCapLine(expectedCapabilities []string, currentCapabilities []string) (errors []error) {
	seen := map[string]struct{}{}
	for _, currentCapability := range currentCapabilities {
		if contains(expectedCapabilities, currentCapability) {
			seen[currentCapability] = struct{}{}
		} else {
			errors = append(errors, &CapabilityNotProvided{Capability: currentCapability})
		}
	}

	for _, expectedCapability := range expectedCapabilities {
		if _, ok := seen[expectedCapability]; !ok {
			errors = append(errors, &UnnecessaryCapability{Capability: expectedCapability})
		}
	}

	return errors
}

func contains(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}
