package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadFileCmd représente la commande upload-file
var uploadFileCmd = &cobra.Command{
	Use:   "upload-file",
	Short: "Uploads a file to a specified S3 bucket via the API",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Usage: upload-file <bucket-name> <file-path>")
		}

		bucketName := args[0]
		filePath := args[1]

		// Récupérer l'URL de l'API via Viper
		apiURL := viper.GetString("s3.api_url")
		if apiURL == "" {
			log.Fatal("API URL is not configured. Please set it in the config file or environment variables.")
		}

		// Lire le fichier à uploader
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()

		// Lire le contenu du fichier
		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		// Extraire le nom du fichier depuis le chemin
		fileName := filepath.Base(filePath)

		// URL de l'API pour l'upload du fichier
		url := fmt.Sprintf("%s/%s/%s", apiURL, bucketName, fileName)

		// Calculer la longueur du contenu
		contentLength := len(fileContent)

		// Faire une requête PUT avec le fichier
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(fileContent))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		// Spécifier le type de contenu du fichier
		req.Header.Set("Content-Type", "application/octet-stream")

		// Ajouter l'en-tête X-Amz-Decoded-Content-Length
		req.Header.Set("X-Amz-Decoded-Content-Length", fmt.Sprintf("%d", contentLength))

		// Envoyer la requête
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error uploading file: %v", err)
		}
		defer resp.Body.Close()

		// Lire le corps de la réponse pour les détails supplémentaires
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		// Vérifier le statut de la réponse
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("File '%s' uploaded successfully to bucket '%s'.\n", fileName, bucketName)
		} else {
			fmt.Printf("Failed to upload file. Status code: %d\n", resp.StatusCode)
			fmt.Printf("Response body: %s\n", string(body))
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadFileCmd)
}
