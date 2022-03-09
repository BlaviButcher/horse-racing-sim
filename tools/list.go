package tools

import (
	"fmt"

	"github.com/blavi/horse/simulation/horsedb"
	"github.com/spf13/cobra"
)

type listCLIOpts struct {
	DatabaseDir string
}

var list = &cobra.Command{
	Use:   "list [db-dir]",
	Short: "Lists all horses in db",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE:  List(),
}

// type MyFunc func(cmd *cobra.Command, args []string) error

func List() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		opts := listCLIOpts{}
		opts.ParseArgs(args)

		database, err := horsedb.OpenDatabase(opts.DatabaseDir)
		if err != nil {
			return fmt.Errorf("opening database: %w", err)
		}
		defer database.Close()

		horses, err := database.ListHorses()
		if err != nil {
			return fmt.Errorf("listing horses from database: %w", err)
		}

		for _, h := range horses {
			fmt.Printf(`
			ID:%s
			Name:%s
			MMR:%d
			`, h.ID, h.Name, h.MMR)
		}
		return nil
	}
}

func (opts *listCLIOpts) ParseArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, got %d", len(args))
	}
	opts.DatabaseDir = args[0]
	return nil
}

func init() {
	rootCmd.AddCommand(list)
}
