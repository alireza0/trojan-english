package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"trojan/core"
	"trojan/trojan"
	"trojan/util"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "trojan",
	Run: func(cmd *cobra.Command, args []string) {
		mainMenu()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check() {
	if !util.IsExists("/usr/local/etc/trojan/config.json") {
		fmt.Println("This machine has no trojan installed! Automatic installation in progress...")
		trojan.InstallTrojan("")
		core.WritePassword(nil)
		trojan.InstallTls()
		trojan.InstallMysql()
		util.SystemctlRestart("trojan-web")
	}
}

func mainMenu() {
	check()
exit:
	for {
		fmt.Println()
		fmt.Println(util.Cyan("Welcome to trojan management GUI"))
		fmt.Println()
		menuList := []string{"trojan Management", "User Management", "Installation Management", "Web Management", "View configuration", "Generate json"}
		switch util.LoopInput("Please choose: ", menuList, false) {
		case 1:
			trojan.ControllMenu()
		case 2:
			trojan.UserMenu()
		case 3:
			trojan.InstallMenu()
		case 4:
			trojan.WebMenu()
		case 5:
			trojan.UserList()
		case 6:
			trojan.GenClientJson()
		default:
			break exit
		}
	}
}
