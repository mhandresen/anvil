package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

func parseMemory(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	mult := int64(1)
	u := strings.ToUpper(s)
	switch {
	case strings.HasSuffix(u, "K"):
		mult = 1024
		u = strings.TrimSuffix(u, "K")
	case strings.HasSuffix(u, "M"):
		mult = 1024 * 1024
		u = strings.TrimSuffix(u, "M")
	case strings.HasSuffix(u, "G"):
		mult = 1024 * 1024 * 1024
		u = strings.TrimSuffix(u, "G")
	}
	n, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory %q", s)
	}
	return n * mult, nil
}