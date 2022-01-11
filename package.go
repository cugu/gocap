package main

import (
	"encoding/json"
	"errors"
	"os/exec"
)

type PackagePublic struct {
	Name       string
	ImportPath string
	Imports    []string
}

func parsePackage(pkgPath string) (imports []string, err error) {
	out, err := exec.Command("go", "list", "--json", pkgPath).Output()
	if err != nil {
		return nil, errors.New("could not parse package")
	}

	pkg := &PackagePublic{}
	if err := json.Unmarshal(out, pkg); err != nil {
		return nil, err
	}
	return pkg.Imports, nil
}

func parsePackagePath(pkgPath string) (string, error) {
	out, err := exec.Command("go", "list", "--json", pkgPath).Output()
	if err != nil {
		return "", errors.New("could not parse package path")
	}

	pkg := &PackagePublic{}
	if err := json.Unmarshal(out, pkg); err != nil {
		return "", err
	}
	return pkg.ImportPath, nil
}
