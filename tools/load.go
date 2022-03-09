package tools

import (
	"fmt"

	"github.com/blavi/horse/simulation"
	"github.com/blavi/horse/simulation/horsedb"
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
	RunE: Load(),
}

func Load() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts := loadCLIOpts{}
		opts.ParseArgs(args)

		database, err := horsedb.OpenDatabase(opts.DatabaseDir)
		if err != nil {
			return fmt.Errorf("opening database: %w", err)
		}
		defer database.Close()

		horses, err := opts.GetRecords()
		if err != nil {
			return fmt.Errorf("getting horses %w", err)
		}

		err = database.SetHorses(horses...)
		if err != nil {
			return fmt.Errorf("setting horses: %w", err)
		}

		horses, _ = simulation.SimulateFreshHorses(horses, 30000)

		totalMMR := 0
		for _, h := range horses {
			fmt.Printf("Name: %s\n", h.Name)
			fmt.Printf("MMR: %d\n", h.MMR)
			fmt.Printf("RawMMR: %d\n", h.RawMMR)
			fmt.Printf("Avg MMR: %d\n", h.AvgMMR)
			totalMMR += h.AvgMMR
		}
		fmt.Printf(`Mean: %d`, totalMMR/len(horses))

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
