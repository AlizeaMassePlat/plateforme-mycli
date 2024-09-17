package cmd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Bucket représente un bucket dans la réponse XML
type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

// ListAllMyBucketsResult représente la structure XML pour lister les buckets
type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

// listBucketsCmd représente la commande `list-buckets`
var listBucketsCmd = &cobra.Command{
	Use:   "list-buckets",
	Short: "List all S3 buckets via the API",
	Run: func(cmd *cobra.Command, args []string) {
		// Récupérer l'URL de l'API depuis le fichier de configuration ou les variables d'environnement
		apiURL := viper.GetString("s3.api_url")
		if apiURL == "" {
			handleError(errors.New("API URL is not configured. Please set it in the config file or environment variables"))
			return
		}

		// Faire une requête GET pour obtenir la liste des buckets
		resp, err := http.Get(apiURL)
		if err != nil {
			handleError(fmt.Errorf("failed to connect to S3 API at %s: %v", apiURL, err))
			return
		}
		defer resp.Body.Close()

		// Vérifier le statut de la réponse
		if resp.StatusCode != http.StatusOK {
			handleError(fmt.Errorf("S3 API returned an error: status code %d", resp.StatusCode))
			return
		}

		// Lire la réponse
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			handleError(fmt.Errorf("failed to read response body: %v", err))
			return
		}

		// Parse la réponse XML
		var result ListAllMyBucketsResult
		err = xml.Unmarshal(body, &result)
		if err != nil {
			handleError(fmt.Errorf("failed to parse XML response: %v", err))
			return
		}

		// Afficher les buckets de manière lisible
		if len(result.Buckets) == 0 {
			fmt.Println("No buckets found.")
			return
		}

		fmt.Println("Buckets:")
		for _, bucket := range result.Buckets {
			fmt.Printf("- %s (created on %s)\n", bucket.Name, bucket.CreationDate)
		}
	},
}

// handleError affiche un message d'erreur et continue l'exécution du programme sans l'arrêter brutalement
func handleError(err error) {
	fmt.Fprintf(log.Writer(), "Error: %v\n", err)
}

func init() {
	// Ajouter la commande `list-buckets` à la racine
	rootCmd.AddCommand(listBucketsCmd)
}
