package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
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
		fmt.Printf("Logpipe - Lightweight Log Monitoring Tool \n")
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

		// Create a separate process group.
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

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

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(signalChan)

		shutdownDone := make(chan struct{})

		go func() {
			select {
			case receivedSignal := <-signalChan:
				fmt.Printf("\n[LOGPIPE] Received signal: %s\n", receivedSignal)

				sig, ok := receivedSignal.(syscall.Signal)
				if !ok {
					sig = syscall.SIGTERM
				}

				// Forward signal to the entire process group.
				if err := syscall.Kill(-cmd.Process.Pid, sig); err != nil {
					fmt.Fprintf(os.Stderr, "[LOGPIPE] Signal error: %v\n", err)
				}

				fmt.Println("[LOGPIPE] Waiting for process shutdown...")

				timer := time.NewTimer(5 * time.Second)
				defer timer.Stop()

				select {
				case <-shutdownDone:
					return

				case <-timer.C:
					fmt.Println("[LOGPIPE] Shutdown timeout. Force killing process...")
					_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				}

			case <-shutdownDone:
				return
			}
		}()

		wg.Wait()

		close(shutdownDone)

		err = cmd.Wait()

		duration := time.Since(startTime)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] Process failed: %v\n", err)

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
