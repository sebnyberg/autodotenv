package autodotenv

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	ErrInvalidRow = errors.New("invalid row")
)

// LoadDotenvIfExists reads environment variables from .env
// and sets the environment using os.Setenv.
//
// If .env contains invalid entries, an error is returned.
// If .env could not be opened, no error is returned.
func LoadDotenv() error {
	f, err := os.OpenFile(".env", os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close .env file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(f)
	n := 0
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.SplitN(row, "=", 2)
		if len(parts) == 1 {
			if parts[0] != "" {
				return fmt.Errorf("%w: failed to parse row %d, value: %v", ErrInvalidRow, n+1, row)
			}
			continue
		}
		if len(os.Getenv(parts[0])) > 0 {
			continue
		}
		os.Setenv(parts[0], parts[1])
		n++
	}
	return nil
}
