package main

import (
	"flag"
	"fmt"
	"os"
	"test/starkbank/mocked/cmd/common"
)

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "Show help")
	flag.Parse()

	if len(os.Args) == 1 {
		errorC()
		return
	}

	if help || (len(os.Args) > 1 && os.Args[1] == "help") {
		helpC()
		return
	}

	command(os.Args[1])
}

func errorC() {
	fmt.Println(common.Red, "Command not specified or not available.", common.Reset)
	fmt.Println(common.Yellow, " Use ./gomd help or ./gomd -h for help.", common.Reset)
}

func helpC() {
	fmt.Println(common.Green, "Commands:")
	fmt.Println("  create:migration")
	fmt.Println("  create:controller")
	fmt.Println("  migrate")
	fmt.Println("     Usage:")
	fmt.Println("       ./gomd create:migration -name <name>")
	fmt.Println("       ./gomd create:controller -name <name>")
	fmt.Println("       ./gomd migrate")
	fmt.Println(common.Reset)
}
