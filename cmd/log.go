package cmd

import (
	"trojan/util"

	"github.com/spf13/cobra"
)

var line int

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View Trojan log",
	Run: func(cmd *cobra.Command, args []string) {
		util.Log("trojan", line)
	},
}

func init() {
	logCmd.Flags().IntVarP(&line, "line", "n", 300, "number of logs to view")
	rootCmd.AddCommand(logCmd)
}
