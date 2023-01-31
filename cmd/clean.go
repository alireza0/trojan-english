package cmd

import (
	"github.com/spf13/cobra"
	"trojan/trojan"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clear designated user traffic",
	Long: `Pass to the designated username to clear the user traffic, and separate multiple user name spaces, such as:
trojan clean zhangsan lisi
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		trojan.CleanDataByName(args)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
