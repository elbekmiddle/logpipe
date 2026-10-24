package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("LogPipe")
	fmt.Println("Local developer log management tool")

	return nil
}
