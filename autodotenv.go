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
	ErrLoadFailed = errors.New("load failed")
	ErrInvalidRow = errors.New("invalid row")
)

// LoadDotenvIfExists reads environment variables from .env
// and sets the environment using os.Setenv. Beware that these variables
// will be passed to child processes.
//
// If .env cannot be loaded, the function returns no error and zero
// loaded variables. This is because .env is expected to only be present
// in the local environment.
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
