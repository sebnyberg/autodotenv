package autodotenv

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var (
	ErrLoadFailed = errors.New("load failed")
	ErrInvalidRow = errors.New("invalid row")
)

// LoadDotenvIfExists reads environment variables from .env
// and sets the environment using os.Setenv. Beware that these variables
// will be passed onto child processes.
// The function returns the number of variables read from the file.
// If no .env file exists, autodotenv.ErrLoadFailed is returned.
func LoadDotenv(fp string) (n int, err error) {
	fp = path.Clean(fp)
	f, err := os.OpenFile(fp, os.O_RDONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("%w: dotenv file at '%v', err: %v", ErrLoadFailed, fp, err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close opened .env file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(f)
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
