package collector

import "time"

type Entry struct {
	Timestamp time.Time
	Stream    string
	Message   string
}
