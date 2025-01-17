package cmd

import (
	"github.com/spf13/cobra"
	"trojan/trojan"
)

// tlsCmd represents the tls command
var tlsCmd = &cobra.Command{
	Use:   "tls",
	Short: "Certificate installation",
	Run: func(cmd *cobra.Command, args []string) {
		trojan.InstallTls()
	},
}

func init() {
	rootCmd.AddCommand(tlsCmd)
}
