package main

import (
	"log"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Generate struct {
		Path string `arg:"" name:"path" help:"Path to the package to generate the go.cap file for"`
	} `cmd:"" help:"Generate go.cap"`
	Check struct {
		Path   string `arg:"" name:"path" help:"Path to the package to check"`
		Ignore string `help:"Ignore the package itself."`
	} `cmd:"" help:"Check go.cap"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "generate <path>":
		generate(CLI.Generate.Path)
	case "check <path>":
		checkCmd(CLI.Check.Path, CLI.Check.Ignore)
	default:
		panic(ctx.Command())
	}
}
