package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// downloadFileCmd représente la commande download-file
var downloadFileCmd = &cobra.Command{
	Use:   "download-file",
	Short: "Downloads a file from a specified S3 bucket via the API",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			log.Fatal("Usage: download-file <bucket-name> <file-name> <destination-path>")
		}

		bucketName := args[0]
		fileName := args[1]
		destPath := args[2]

		// Construire l'URL pour l'API
		url := fmt.Sprintf("http://localhost:9090/%s/%s", bucketName, fileName)

		// Effectuer une requête GET
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error making request to download file: %v", err)
		}
		defer resp.Body.Close()

		// Vérifier le code de statut
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to download file. Status code: %d", resp.StatusCode)
		}

		// Ouvrir le fichier de destination
		out, err := os.Create(filepath.Join(destPath, fileName))
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer out.Close()

		// Copier le contenu de la réponse dans le fichier
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalf("Error writing file: %v", err)
		}

		fmt.Printf("File '%s' downloaded successfully to '%s'.\n", fileName, destPath)
	},
}

func init() {
	rootCmd.AddCommand(downloadFileCmd)
}
