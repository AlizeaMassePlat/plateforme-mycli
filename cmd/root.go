package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "my-cli",
	Short: "My CLI tool",
	Long: `My CLI tool is a command line interface that allows users
to manage various tasks efficiently.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from My CLI!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
