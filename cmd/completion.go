package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Automatically command to make up (support BASH and ZSH)",
	Long: `
Command support for bash and zsh
a. bash: Add the following command to ~/.bashrc 
   source <(trojan completion bash)

b. zsh: Add the following command to ~/.zshrc
   source <(trojan completion zsh)
`,
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(&cobra.Command{Use: "bash", Short: "bash command to make up", Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	}})
	completionCmd.AddCommand(&cobra.Command{Use: "zsh", Short: "zsh command to make up", Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
	}})
}
