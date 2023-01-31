package cmd

import (
	"trojan/trojan"

	"github.com/spf13/cobra"
)

// portCmd represents the info command
var portCmd = &cobra.Command{
	Use:   "port",
	Short: "Modify the trojan port",
	Run: func(cmd *cobra.Command, args []string) {
		trojan.ChangePort()
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
}
