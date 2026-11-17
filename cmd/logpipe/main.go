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
		command, args, logPath, err := parseRunArgs(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] %v\n", err)
			os.Exit(1)
		}

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

func parseRunArgs(args []string) (string, []string, string, error) {
	logPath, err := defaultLogPath()
	if err != nil {
		return "", nil, "", err
	}

	if len(args) > 0 && args[0] == "--log-file" {
		if len(args) < 2 || args[1] == "--" {
			return "", nil, "", fmt.Errorf("--log-file requires a file path")
		}

		logPath = args[1]
		args = args[2:]
	}

	if len(args) < 2 || args[0] != "--" {
		return "", nil, "", fmt.Errorf("usage: logpipe run [--log-file <path>] -- <command> [args...]")
	}

	command := args[1]
	commandArgs := args[2:]

	return command, commandArgs, logPath, nil
}
