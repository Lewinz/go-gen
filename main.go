package main

import (
	"fmt"
	"os"

	"github.com/lewinz/go-gen/model"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "go-gen",
	Short: "A code generation tool for Go projects",
	Long: `A flexible code generation tool that supports multiple generators.
It can generate code for different purposes like models, services, and more.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\ncommit: %s\nbuilt at: %s\n", version, commit, date)
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(model.GetModelCmd())
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
