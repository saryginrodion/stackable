package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	if d.Nanoseconds() < 1000 {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}

	if d.Microseconds() < 1000 {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}

	if d.Milliseconds() < 1000 {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}

	if d.Seconds() < 60 {
		return fmt.Sprintf("%fs", d.Seconds())
	}

	if d.Minutes() < 60 {
		return fmt.Sprintf("%fs", d.Minutes())
	}

	return fmt.Sprintf("%fh", d.Hours())
}
