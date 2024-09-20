package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// helloCmd représente la commande hello
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Prints a hello message",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func init() {
	RootCmd.AddCommand(helloCmd)
}
