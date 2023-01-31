package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"trojan/core"
	"trojan/trojan"
)

// upgradeCmd represents the update command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade database and Trojan configuration file",
}

func upgradeConfig() {
	domain, _ := core.GetValue("domain")
	if domain == "" {
		return
	}
	core.WriteDomain(domain)
	trojan.Restart()
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.AddCommand(&cobra.Command{Use: "db", Short: "Upgrade database", Run: func(cmd *cobra.Command, args []string) {
		if err := core.GetMysql().UpgradeDB(); err != nil {
			fmt.Println(err)
		}
	}})
	upgradeCmd.AddCommand(&cobra.Command{Use: "config", Short: "Upgrade configuration file", Run: func(cmd *cobra.Command, args []string) {
		upgradeConfig()
	}})
}
