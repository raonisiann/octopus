package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting with Go")

	if len(os.Args) < 2 {
		fmt.Printf("Show help")
		os.Exit(-1)
	}

	switch cmd := os.Args[1]; cmd {
	case "apply":
		fmt.Println("Apply sub-command")
		ApplyCmd(os.Args[2:])
	case "help":
		fmt.Println("Show help")
	case "validate":
		fmt.Println("Validate sub-command")
	default:
		fmt.Printf("Unknown option '%s'\n", cmd)
	}

}
