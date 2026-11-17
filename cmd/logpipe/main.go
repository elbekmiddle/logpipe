package main

import (
	"fmt"
	"os"
	"path/filepath"

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
		fmt.Println("Logpipe - Lightweight Log Monitoring Tool")
		fmt.Println("Commands:")
		fmt.Println("  --version, -v    Show version")
		fmt.Println("  --help, -h       Show help")
		fmt.Println("  run              Run backend process")

	case "run":
		runArgs := os.Args[2:]

		logPath, err := defaultLogPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] %v\n", err)
			os.Exit(1)
		}

		if len(runArgs) >= 2 && runArgs[0] == "--log-file" {
			logPath = runArgs[1]
			runArgs = runArgs[2:]
		}

		if len(runArgs) < 2 || runArgs[0] != "--" {
			fmt.Println("Usage: logpipe run [--log-file <path>] -- <command> [args...]")
			return
		}

		command := runArgs[1]
		args := runArgs[2:]

		if err := process.Run(command, args, logPath); err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}

func defaultLogPath() (string, error) {
	stateDir := os.Getenv("XDG_STATE_HOME")

	if stateDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home directory: %w", err)
		}

		stateDir = filepath.Join(homeDir, ".local", "state")
	}

	return filepath.Join(stateDir, "logpipe", "logs.jsonl"), nil
}
