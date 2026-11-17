package collector

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

func Stream(reader io.Reader, output io.Writer, stream string) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		entry := Entry{
			Timestamp: time.Now(),
			Stream:    stream,
			Message:   scanner.Text(),
		}

		if _, err := fmt.Fprintf(
			output,
			"[%s] %s\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Stream,
			entry.Message,
		); err != nil {
			return fmt.Errorf("write %s: %w", stream, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read %s: %w", stream, err)
	}

	return nil
}
