package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadFileCmd représente la commande upload-file
var UploadFileCmd = &cobra.Command{
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

		// Extraire le nom du fichier depuis le chemin
		fileName := filepath.Base(filePath)

		// URL de l'API pour l'upload du fichier
		url := fmt.Sprintf("%s/%s/%s", apiURL, bucketName, fileName)

		// Obtenir la taille du fichier pour la barre de progression
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatalf("Error getting file info: %v", err)
		}
		totalSize := fileInfo.Size()

		// Créer le client HTTP
		client := &http.Client{}
		pipeReader, pipeWriter := io.Pipe()

		// Préparer la requête HTTP
		req, err := http.NewRequest("PUT", url, pipeReader)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		req.ContentLength = totalSize
		req.Header.Set("X-Amz-Decoded-Content-Length", fmt.Sprintf("%d", totalSize))

		// Lancer l'upload et mettre à jour la barre de progression
		progressChan := make(chan struct{})
		var uploadedSize int64 = 0

		// Fonction pour démarrer la barre de progression
		go func() {
			for {
				select {
				case <-progressChan:
					return
				default:
					printProgressUpload(uploadedSize, totalSize)
					time.Sleep(100 * time.Millisecond) // Rafraîchir toutes les 100ms
				}
			}
		}()

		// Lire et écrire le fichier en streaming
		go func() {
			defer pipeWriter.Close() // Fermer l'écriture après avoir fini l'upload
			buffer := make([]byte, 256)
			for {
				n, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					progressChan <- struct{}{} // Arrêter l'animation
					pipeWriter.CloseWithError(fmt.Errorf("failed to read file: %w", err))
					return
				}
				if n == 0 {
					break
				}
				_, err = pipeWriter.Write(buffer[:n])
				if err != nil {
					progressChan <- struct{}{} // Arrêter l'animation
					pipeWriter.CloseWithError(fmt.Errorf("failed to write to pipe: %w", err))
					return
				}
				uploadedSize += int64(n)
			}
		}()

		// Envoyer la requête
		resp, err := client.Do(req)
		if err != nil {
			progressChan <- struct{}{} // Arrêter l'animation
			log.Fatalf("Error uploading file: %v", err)
		}
		defer resp.Body.Close()

		// Mettre à jour la barre de progression pour indiquer 100%
		printProgressUpload(uploadedSize, totalSize)
		progressChan <- struct{}{} // Arrêter la barre de progression

		// Vérifier le statut de la réponse et afficher le message après l'upload
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("\nFile '%s' uploaded successfully to bucket '%s'.\n", fileName, bucketName)
		} else {
			fmt.Printf("Failed to upload file. Status code: %d\n", resp.StatusCode)
		}
	},
}

func init() {
	RootCmd.AddCommand(UploadFileCmd)
}

// Fonction pour afficher la barre de progression
func printProgressUpload(uploaded, total int64) {
	if total == -1 {
		fmt.Printf("\rUploading... %d bytes", uploaded)
		return
	}
	percentage := float64(uploaded) / float64(total) * 100
	progressBar := int(percentage / 2) // barre de 50 caractères
	fmt.Printf("\r[%-50s] %3.2f%%", strings.Repeat("#", progressBar), percentage)
}


