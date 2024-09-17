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
)

// downloadFileCmd représente la commande download-file
var downloadFileCmd = &cobra.Command{
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

		// Construire l'URL pour l'API
		url := fmt.Sprintf("http://localhost:9090/%s/%s", bucketName, fileName)

		// Télécharger le fichier
		if err := downloadFile(url, destPath, fileName); err != nil {
			log.Printf("Error downloading file: %v", err)
		} else {
			fmt.Printf("File '%s' downloaded successfully to '%s'.\n", fileName, destPath)
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

    // Vérifier le code de statut HTTP
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to download file. Status code: %d", resp.StatusCode)
    }

    // Ouvrir le fichier de destination
    out, err := os.Create(filepath.Join(destPath, fileName))
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer out.Close()

    // Créer un buffer pour lire le contenu
    buffer := make([]byte, 256) // Taille de buffer réduite

    // Obtenir la taille du fichier pour l'animation de progression
    totalSize := resp.ContentLength
    var downloadedSize int64 = 0

    // Lire et copier le contenu avec une animation de progression
    firstRead := true // Indicateur pour démarrer la barre de progression après la première lecture
    progressChan := make(chan struct{})

    for {
        n, err := resp.Body.Read(buffer)
        if err != nil && err != io.EOF {
            progressChan <- struct{}{} // Arrêter l'animation
            return fmt.Errorf("failed to read response body: %w", err)
        }
        if n == 0 {
            break
        }

        // Écrire dans le fichier de destination
        if _, err := out.Write(buffer[:n]); err != nil {
            progressChan <- struct{}{} // Arrêter l'animation
            return fmt.Errorf("failed to write to file: %w", err)
        }
        downloadedSize += int64(n)

        // Démarrer l'animation de progression après la première lecture
        if firstRead {
            firstRead = false
            go func() {
                for {
                    select {
                    case <-progressChan:
                        return
                    default:
                        printProgress(downloadedSize, totalSize)
                        time.Sleep(100 * time.Millisecond) // rafraîchir toutes les 100ms
                    }
                }
            }()
        }
    }

    // Mettre à jour une dernière fois la barre de progression pour indiquer 100%
    printProgress(downloadedSize, totalSize)

    // Arrêter l'animation de progression
    progressChan <- struct{}{}
	fmt.Print("\n")
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
	rootCmd.AddCommand(downloadFileCmd)
}
