package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Object represents an object in the S3 bucket
type Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	Size         int    `xml:"Size"`
}

// ListBucketResult represents the response structure from the S3 API
type ListBucketResult struct {
	Objects []Object `xml:"Contents"`
}

// listObjectCmd represents the list-object command
var ListObjectCmd = &cobra.Command{
	Use:   "list-object",
	Short: "List objects in a specified S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Usage: list-object <bucket-name>")
		}

		bucketName := args[0]

		// Récupérer l'URL de l'API via Viper
		apiURL := viper.GetString("s3.api_url")
		if apiURL == "" {
			log.Fatal("API URL is not configured. Please set it in the config file or environment variables.")
		}

		// Construire l'URL pour lister les objets du bucket
		url := fmt.Sprintf("%s/%s/", apiURL, bucketName)

		// Effectuer une requête GET
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error making request to list objects: %v", err)
		}
		defer resp.Body.Close()

		// Vérifier le code de statut
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to list objects. Status code: %d", resp.StatusCode)
		}

		// Lire le corps de la réponse
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		// Parser la réponse XML
		var result ListBucketResult
		err = xml.Unmarshal(body, &result)
		if err != nil {
			log.Fatalf("Error parsing XML response: %v", err)
		}
		

		// Afficher les objets
		fmt.Println("Objects:")
		const readableDateLayout = "2006-01-02 15:04:05"
		const inputLayout = time.RFC3339 // Le format attendu

		for _, obj := range result.Objects {
			// Convertir la chaîne de caractères LastModified en time.Time
			lastModifiedTime, err := time.Parse(inputLayout, obj.LastModified)
			if err != nil {
				log.Printf("Error parsing date: %v", err)
				continue
			}

			// Formater la date pour un affichage lisible
			readableDate := lastModifiedTime.Format(readableDateLayout)
			fmt.Printf("- [%s] %dB %s\n", readableDate, obj.Size, obj.Key)
		}
	},
}

func init() {
	RootCmd.AddCommand(ListObjectCmd)
}
