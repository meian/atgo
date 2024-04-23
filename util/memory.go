package util

import "fmt"

func ParseMemory(s string) (int, error) {
	var n int
	if _, err := fmt.Sscanf(s, "%d KB", &n); err == nil {
		return n << 10, nil
	}
	if _, err := fmt.Sscanf(s, "%d MB", &n); err == nil {
		return n << 20, nil
	}
	if _, err := fmt.Sscanf(s, "%d GB", &n); err == nil {
		return n << 30, nil
	}
	if _, err := fmt.Sscanf(s, "%d B", &n); err == nil {
		return n, nil
	}
	return 0, fmt.Errorf("failed to parse memory: %s", s)
}

func FormatMemory(memBytes int) string {
	if memBytes <= 1<<10 {
		return fmt.Sprintf("%d B", memBytes)
	}
	if memBytes <= 1<<20 {
		return fmt.Sprintf("%d KB", memBytes>>10)
	}
	if memBytes <= 1<<30 {
		return fmt.Sprintf("%d MB", memBytes>>20)
	}
	return fmt.Sprintf("%d GB", memBytes>>30)
}
