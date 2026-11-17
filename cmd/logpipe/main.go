package main

import (
	"fmt"
	"os"
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

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}