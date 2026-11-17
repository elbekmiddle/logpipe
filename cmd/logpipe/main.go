package main

import (
	"fmt"
	"os"

	"logpipe/internal/process"
)

const VERSION = "0.0.1"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Logpipe - Lightweight Log Monitoring Tool")
		fmt.Println("Usage: logpipe [command]")
		return
	}

	switch os.Args[1] {
	case "--version", "-v":
		fmt.Printf("Logpipe version %s\n", VERSION)

	case "--help", "-h":
		fmt.Println("Logpipe - Lightweight Log Monitoring Tool \n")
		fmt.Println("Commands:")
		fmt.Println("  --version, -v    Show version")
		fmt.Println("  --help, -h       Show help")
		fmt.Println("  run              Run backend process")

	case "run":
		runArgs := os.Args[2:]

		if len(runArgs) < 2 || runArgs[0] != "--" {
			fmt.Println("Usage: logpipe run -- <command> [args...]")
			return
		}

		command := runArgs[1]
		args := runArgs[2:]

		if err := process.Run(command, args); err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}
