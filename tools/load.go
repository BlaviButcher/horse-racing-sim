package tools

import (
	"fmt"

	"github.com/blavi/horse/horse/horsedb"
	"github.com/spf13/cobra"
)

type loadCLIOpts struct {
	DatabaseDir string
	InputFile   string
}

var load = &cobra.Command{
	Use:   "load [db-dir] [file-path]",
	Short: "Loads a list of horses",
	Long: `db-dir is directory that database is found in
	File must be of *.txt and contain a single name for each line.
	The program will automatically supply and ID and MMR rating`,
	Args: cobra.MinimumNArgs(1),
	RunE: Run(),
}

func Run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts := loadCLIOpts{}
		opts.ParseArgs(args)

		database, err := horsedb.OpenDatabase(opts.DatabaseDir)
		if err != nil {
			return fmt.Errorf("opening database: %w", err)
		}
		defer database.Close()

		opts.GetRecords()

		return nil

	}
}

// ParseArgs checks the required args are present
func (opts *loadCLIOpts) ParseArgs(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", len(args))
	}
	opts.DatabaseDir = args[0]
	opts.InputFile = args[1]
	return nil
}

func init() {
	rootCmd.AddCommand(load)
}
