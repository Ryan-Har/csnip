package common

import (
	"fmt"
	"os"
)

func ReadFromFile(s string) (string, error) {
	data, err := os.ReadFile(s)
	if err != nil {
		return "", fmt.Errorf("unable to read file: %w", err)
	}

	return string(data), nil
}
