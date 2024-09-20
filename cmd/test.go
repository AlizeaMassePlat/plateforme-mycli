package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// testCmd représente la commande `test`
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Execute the unit tests for the CLI",
	Long:  `This command runs all the unit tests defined in the project using the 'go test' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Si des arguments sont fournis, exécute uniquement ces tests
		if len(args) > 0 {
			for _, testName := range args {
				err := runSpecificTest(testName)
				if err != nil {
					fmt.Printf("Error running test %s: %v\n", testName, err)
					os.Exit(1)
				}
			}
			fmt.Println("Specified tests passed successfully!")
		} else {
			// Si aucun argument n'est fourni, exécute tous les tests
			err := runTests()
			if err != nil {
				fmt.Printf("Error running tests: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("All tests passed successfully!")
		}
	},
}

// runSpecificTest exécute un test spécifique
func runSpecificTest(testName string) error {
	// Exécuter un test spécifique en utilisant `go test -run <TestName>`
	cmd := exec.Command("go", "test", "-run", testName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(testCmd)
}

// runTests exécute la commande `go test ./...` pour exécuter tous les tests du projet
func runTests() error {
	// Préparer la commande `go test`
	cmd := exec.Command("go", "test", "./...")

	// Rediriger la sortie standard et les erreurs de la commande vers la sortie de la CLI
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Exécuter la commande `go test` et retourner le résultat
	return cmd.Run()
}
