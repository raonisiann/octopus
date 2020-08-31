package main

import (
	"fmt"
	"os"
	"parser"
)

// ApplyCmd get the arguments from apply command
func ApplyCmd(args []string) {

	if len(args) == 0 {
		fmt.Println("File required")
		os.Exit(-1)
	}

	fileName := args[0]
	parser.Parse(fileName)
}
