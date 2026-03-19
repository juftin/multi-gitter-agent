package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multi-gitter-agent",
	Short: "Run AI prompts across multiple repositories using multi-gitter",
}

var runCmd = &cobra.Command{
	Use:   "run [prompt]",
	Short: "Run a prompt across repositories",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		fmt.Printf("Running prompt: %s\n", prompt)
		// TODO: Implement multi-gitter and AI agent integration
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	// Additional flags like --org, --repo, etc. will go here.
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
