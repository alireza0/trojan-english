package cmd

import (
	"github.com/spf13/cobra"
	"trojan/trojan"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop trojan",
	Run: func(cmd *cobra.Command, args []string) {
		trojan.Stop()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
