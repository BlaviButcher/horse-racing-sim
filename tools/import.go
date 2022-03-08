package tools

import "github.com/spf13/cobra"

var importCmd = &cobra.Command{
	Use:   "import [file-path]",
	Short: "Import a list of horses",
	Long: `file must be of *.txt and contain a single name for each line.
	The program will automatically supply and ID and MMR rating`,
	Args: cobra.MinimumNArgs(1),
	RunE: Run(),
}

func Run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
