package cmd

import (
	"trojan/trojan"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "View trojan status",
	Run: func(cmd *cobra.Command, args []string) {
		trojan.Status(true)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
