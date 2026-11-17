package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
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

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] stdout error: %v\n", err)
			os.Exit(1)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] stderr error: %v\n", err)
			os.Exit(1)
		}

		startTime := time.Now()

		fmt.Printf("[LOGPIPE] Starting process: %s\n", command)

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] Failed to start: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("[LOGPIPE] Process PID: %d\n", cmd.Process.Pid)

		var wg sync.WaitGroup

		wg.Add(2)

		go streamLogs(stdout, os.Stdout, "STDOUT", &wg)
		go streamLogs(stderr, os.Stderr, "STDERR", &wg)

		err = cmd.Wait()

		wg.Wait()

		duration := time.Since(startTime)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] Process failed: %v\n", err)
		    fmt.Printf("[LOGPIPE] Duration: %s\n", duration.Round(time.Millisecond))

			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Printf("[LOGPIPE] Exit code: %d\n", exitErr.ExitCode())
			}

			os.Exit(1)
		}

		fmt.Println("[LOGPIPE] Process finished successfully")
		fmt.Printf("[LOGPIPE] Duration: %s\n", duration.Round(time.Millisecond))
		fmt.Println("[LOGPIPE] Exit code: 0")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}

func streamLogs(reader io.Reader, output io.Writer, stream string, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Fprintf(output, "[%s] %s\n", stream, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "[LOGPIPE] %s read error: %v\n", stream, err)
	}
}