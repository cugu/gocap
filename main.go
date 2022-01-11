package main

import (
	"log"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Generate struct {
		Path string `arg:"" name:"path" help:"..."`
	} `cmd:"" help:"Generate go.cap"`
	Check struct {
		Path string `arg:"" name:"path" help:"..."`
	} `cmd:"" help:"Check go.cap"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "generate <path>":
		generate(CLI.Generate.Path)
	case "check <path>":
		checkCmd(CLI.Check.Path)
	default:
		panic(ctx.Command())
	}
}
