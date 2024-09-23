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

// downloadFileCmd représente la commande download-file
var DownloadFileCmd = &cobra.Command{
	Use:   "download-file",
	Short: "Downloads a file from a specified S3 bucket via the API",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			log.Println("Usage: download-file <bucket-name> <file-name> <destination-path>")
			return
		}

		bucketName := args[0]
		fileName := args[1]
		destPath := args[2]

		// Récupérer l'URL de l'API depuis la configuration via Viper
		apiURL := viper.GetString("s3.api_url")
		if apiURL == "" {
			log.Fatal("API URL is not configured. Please set it in the config file or environment variables.")
		}

		// Construire l'URL pour l'API
		url := fmt.Sprintf("%s/%s/%s", apiURL, bucketName, fileName)

		// Télécharger le fichier
		err := downloadFile(url, destPath, fileName)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	},
}

func downloadFile(url, destPath, fileName string) error {
	// Faire la requête HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Gérer les statuts HTTP
	switch resp.StatusCode {
	case http.StatusOK:
		// Statut 200 OK - Procéder au téléchargement
		fmt.Printf("File '%s' is being downloaded...\n", fileName)

		// Ouvrir le fichier de destination
		out, err := os.Create(filepath.Join(destPath, fileName))
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer out.Close()

		// Créer un buffer pour lire le contenu
		buffer := make([]byte, 256) // Taille de buffer
		totalSize := resp.ContentLength
		var downloadedSize int64 = 0

		// Canal pour gérer la progression
		progressChan := make(chan struct{})

		// Goroutine pour afficher la barre de progression
		go func() {
			for {
				select {
				case <-progressChan:
					return
				default:
					printProgress(downloadedSize, totalSize)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()

		// Lire et copier le contenu
		for {
			n, err := resp.Body.Read(buffer)
			if err != nil && err != io.EOF {
				progressChan <- struct{}{} // Arrêter la progression
				return fmt.Errorf("failed to read response body: %w", err)
			}
			if n == 0 {
				break
			}

			// Écrire dans le fichier de destination
			if _, err := out.Write(buffer[:n]); err != nil {
				progressChan <- struct{}{} // Arrêter la progression
				return fmt.Errorf("failed to write to file: %w", err)
			}
			downloadedSize += int64(n)
		}

		// Mettre à jour une dernière fois la barre de progression
		printProgress(downloadedSize, totalSize)
		progressChan <- struct{}{} // Arrêter l'animation de progression
		fmt.Println("\nDownload completed successfully.")

	case http.StatusNotFound:
		// Statut 404 Not Found - Fichier introuvable
		return fmt.Errorf("the system cannot find the file specified (404 Not Found)")

	case http.StatusInternalServerError:
		// Statut 500 Internal Server Error
		return fmt.Errorf("internal server error (500)")

	default:
		// Tout autre statut
		return fmt.Errorf("failed to download file. Status code: %d", resp.StatusCode)
	}

	return nil
}

// Affiche une barre de progression simple
func printProgress(downloaded, total int64) {
	if total == -1 {
		fmt.Printf("\rDownloading... %d bytes", downloaded)
		return
	}
	percentage := float64(downloaded) / float64(total) * 100
	progressBar := int(percentage / 2) // barre de 50 caractères
	fmt.Printf("\r[%-50s] %3.2f%%", strings.Repeat("#", progressBar), percentage)
}

func init() {
	RootCmd.AddCommand(DownloadFileCmd)
}
