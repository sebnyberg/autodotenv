package autodotenv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrInvalidRow = errors.New("invalid row")
)

// LoadDotenvIfExists reads environment variables from .env
// and sets values using os.Setenv.
//
// Returns the number of entries loaded and an error marking
// whether the file exists and is malformatted. Missing files
// do not cause any error.
func LoadDotenvIfExists() (int, error) {
	f, err := os.OpenFile(".env", os.O_RDONLY, 0644)
	if err != nil {
		return 0, nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	n := 0
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.SplitN(row, "=", 2)
		if len(parts) == 1 {
			if parts[0] != "" {
				return n, fmt.Errorf("%w: failed to parse row %d, value: %v", ErrInvalidRow, n+1, row)
			}
			continue
		}
		if len(os.Getenv(parts[0])) > 0 {
			continue
		}
		os.Setenv(parts[0], parts[1])
		n++
	}
	return n, nil
}
