package main

import (
	"fmt"
	"os"
	"os/exec"
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
		fmt.Println("Logpipe - Lightweight Log Monitoring Tool")
		fmt.Println()
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

		cmd := exec.Command(command, args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		fmt.Printf("[LOGPIPE] Starting process: %s\n", command)

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] Process failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("[LOGPIPE] Process finished successfully")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}