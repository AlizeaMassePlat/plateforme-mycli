package cmd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
var ListBucketsCmd = &cobra.Command{
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

		// Afficher les buckets
		fmt.Println("Buckets:")
		const readableDateLayout = "2006-01-02 15:04:05"
		const inputLayout = time.RFC3339 // Le format attendu : "2024-09-17T08:58:31Z")
		
		for _, bucket := range result.Buckets {
			// Convertir la chaîne de caractères CreationDate en time.Time
			creationDateTime, err := time.Parse(inputLayout, bucket.CreationDate)
			if err != nil {
				log.Printf("Error parsing date: %v", err)
				continue 
			}
			
			// Formater la date pour un affichage lisible
			readableDate := creationDateTime.Format(readableDateLayout)
			fmt.Printf("- [%s] %s\n", readableDate, bucket.Name)
		}
	},
}

// handleError affiche un message d'erreur et continue l'exécution du programme sans l'arrêter brutalement
func handleError(err error) {
	fmt.Fprintf(log.Writer(), "Error: %v\n", err)
}

func init() {
	// Ajouter la commande `list-buckets` à la racine
	RootCmd.AddCommand(ListBucketsCmd)
}
