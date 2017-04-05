package main

import (
	"fmt"
	"os"

	"github.com/lara-go/installer/commands"
)

const (
	version = "1.0.0"
	help    = `
  LaraGo Installer (v%s)

  Usage: larago <command> [options]

  Commands:
    install    generate a new project from the boilerplate
`
)

var verbose bool

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		printHelp()
	}

	switch argsWithoutProg[0] {
	case "install":
		commands.Install(argsWithoutProg, verbose)
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Printf(help, version)
	os.Exit(0)
}
