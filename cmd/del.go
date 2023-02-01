package cmd

import (
	"trojan/trojan"

	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete user",
	Run: func(cmd *cobra.Command, args []string) {
		trojan.DelUser()
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
