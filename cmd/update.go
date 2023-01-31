package cmd

import (
	"github.com/spf13/cobra"
	"trojan/trojan"
)

// upgradeCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update trojan",
	Long:  "Can be added with a specific version of the update, such as 'Trojan Update v0.10.0', and the latest version is installed without adding a version number",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := ""
		if len(args) == 1 {
			version = args[0]
		}
		trojan.InstallTrojan(version)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
