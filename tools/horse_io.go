package tools

import (
	"bufio"
	"fmt"
	"os"

	"github.com/blavi/horse/simulation/horsedb"
)

// getRecords opens and reads file into an Animal struct.
// Returns an Animal array and error
func (opts *loadCLIOpts) GetRecords() ([]*horsedb.Horse, error) {

	horses := []*horsedb.Horse{}

	file, err := os.Open(opts.InputFile)
	if err != nil {
		return nil, fmt.Errorf("opening file %w", err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		horses = append(horses, horsedb.NewHorse(scanner.Text()))
	}

	return horses, scanner.Err()
}
