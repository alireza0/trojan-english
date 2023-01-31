package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"trojan/util"
	"trojan/web"
)

var (
	host    string
	port    int
	ssl     bool
	timeout int
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start up with web  ",
	Run: func(cmd *cobra.Command, args []string) {
		web.Start(host, port, timeout, ssl)
	},
}

func init() {
	webCmd.Flags().StringVarP(&host, "host", "", "0.0.0.0", "web service monitoring address")
	webCmd.Flags().IntVarP(&port, "port", "p", 80, "web service Start up port")
	webCmd.Flags().BoolVarP(&ssl, "ssl", "", false, "If the web service runs in HTTPS")
	webCmd.Flags().IntVarP(&timeout, "timeout", "t", 120, "Login timeout time (min)")
	webCmd.AddCommand(&cobra.Command{Use: "stop", Short: "Stop trojan-web", Run: func(cmd *cobra.Command, args []string) {
		util.SystemctlStop("trojan-web")
	}})
	webCmd.AddCommand(&cobra.Command{Use: "start", Short: "Start up trojan-web", Run: func(cmd *cobra.Command, args []string) {
		util.SystemctlStart("trojan-web")
	}})
	webCmd.AddCommand(&cobra.Command{Use: "restart", Short: "Restart trojan-web", Run: func(cmd *cobra.Command, args []string) {
		util.SystemctlRestart("trojan-web")
	}})
	webCmd.AddCommand(&cobra.Command{Use: "status", Short: "View Trojan-WEB status", Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(util.SystemctlStatus("trojan-web"))
	}})
	webCmd.AddCommand(&cobra.Command{Use: "log", Short: "View Trojan-WEB logs", Run: func(cmd *cobra.Command, args []string) {
		util.Log("trojan-web", 300)
	}})
	rootCmd.AddCommand(webCmd)
}
