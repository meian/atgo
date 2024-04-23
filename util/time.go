package util

import (
	"fmt"
	"time"
)

func ParseHoursMinutes(s string) (time.Duration, error) {
	var h, m int
	if _, err := fmt.Sscanf(s, "%d:%d", &h, &m); err != nil {
		return time.Duration(0), err
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func ParseDuration(s string) (time.Duration, error) {
	var n int
	if _, err := fmt.Sscanf(s, "%d sec", &n); err == nil {
		return time.Duration(n) * time.Second, nil
	}
	if _, err := fmt.Sscanf(s, "%d msec", &n); err == nil {
		return time.Duration(n) * time.Millisecond, nil
	}
	var nf float64
	if _, err := fmt.Sscanf(s, "%f sec", &nf); err == nil {
		return time.Duration(nf * 1000.0 * float64(time.Millisecond)), nil
	}
	if _, err := fmt.Sscanf(s, "%f msec", &nf); err == nil {
		return time.Duration(nf * 1000.0 * float64(time.Microsecond)), nil
	}
	return time.Duration(0), fmt.Errorf("failed to parse duration: %s", s)
}

func FormatDuration(d time.Duration) string {
	if d.Truncate(time.Second) == d {
		return fmt.Sprintf("%d sec", int(d.Seconds()))
	}
	return fmt.Sprintf("%d msec", int(d.Milliseconds()))
}
