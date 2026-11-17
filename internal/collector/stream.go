package collector

import (
	"bufio"
	"fmt"
	"io"
)

func Stream(reader io.Reader, output io.Writer, stream string) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if _, err := fmt.Fprintf(output, "[%s] %s\n", stream, scanner.Text()); err != nil {
			return fmt.Errorf("write %s: %w", stream, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read %s: %w", stream, err)
	}

	return nil
}
