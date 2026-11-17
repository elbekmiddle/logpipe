package process

import (
	"bufio"
	"fmt"
	"io"
	"logpipe/internal/collector"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run(command string, args []string) error {
	cmd := exec.Command(command, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	startTime := time.Now()

	fmt.Printf("[LOGPIPE] Starting process: %s\n", command)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	fmt.Printf("[LOGPIPE] Process PID: %d\n", cmd.Process.Pid)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		if err := collector.Stream(stdout, os.Stdout, "STDOUT"); err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] %v\n", err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := collector.Stream(stderr, os.Stderr, "STDERR"); err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] %v\n", err)
		}
	}()

	waitChan := make(chan error, 1)

	go func() {
		wg.Wait()
		waitChan <- cmd.Wait()
	}()

	var waitErr error

	select {
	case waitErr = <-waitChan:

	case receivedSignal := <-signalChan:
		fmt.Printf("\n[LOGPIPE] Received signal: %s\n", receivedSignal)

		sig, ok := receivedSignal.(syscall.Signal)
		if !ok {
			sig = syscall.SIGTERM
		}

		if err := syscall.Kill(-cmd.Process.Pid, sig); err != nil {
			fmt.Fprintf(os.Stderr, "[LOGPIPE] Signal error: %v\n", err)
		}

		fmt.Println("[LOGPIPE] Waiting for process shutdown...")

		timer := time.NewTimer(5 * time.Second)

		select {
		case waitErr = <-waitChan:

		case <-timer.C:
			fmt.Println("[LOGPIPE] Shutdown timeout. Force killing process...")

			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

			waitErr = <-waitChan
		}

		timer.Stop()
	}

	duration := time.Since(startTime)

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			fmt.Printf("[LOGPIPE] Exit code: %d\n", exitErr.ExitCode())
		}

		return fmt.Errorf("process failed: %w", waitErr)
	}

	fmt.Println("[LOGPIPE] Process finished successfully")
	fmt.Printf("[LOGPIPE] Duration: %s\n", duration.Round(time.Millisecond))
	fmt.Println("[LOGPIPE] Exit code: 0")

	return nil
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
