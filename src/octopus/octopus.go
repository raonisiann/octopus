package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Starting with Go")

    if len(os.Args) < 2 {
        showGlobalHelp()
        os.Exit(-1)
    }

    switch cmd := os.Args[1]; cmd {
    case "apply":
        fmt.Println("Apply sub-command")
        ApplyCmd(os.Args[2:])
    case "help":
        showGlobalHelp()
    case "validate":
        fmt.Println("Validate sub-command")
    default:
        fmt.Printf("Unknown option '%s'\n", cmd)
        showGlobalHelp()
    }

}

func showGlobalHelp() {
    helpText := `
Octopus vX.X.X

Basic description about what is octopus and what 
can be done with it. Should be no more than 4 lines. 

Online documentation at: https://github.com/raonisiann/octopus

Main Commands:
    apply              Apply resource to target
    config             Configure parameters on agent
    get                Get catalog from Octopus server (requires server setup)
    help               Show this help
    import             Import existing resource into class (requires connectivity with target)
    plan               Generated differences between current state and new state
    pull               Pull catalog from remote Octopus repository
    validate           Validate catalog file
    version            Show agent version

`
    fmt.Println(helpText)
}
