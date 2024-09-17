package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/spf13/cobra"
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
		// URL de votre API pour lister les buckets
		url := "http://localhost:9090/"

		// Faire une requête GET pour obtenir la liste des buckets
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		// Lire la réponse
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		// Vérifier le statut de la réponse
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to list buckets. Status code: %d\n", resp.StatusCode)
			return
		}

		// Parse la réponse XML
		var result ListAllMyBucketsResult
		err = xml.Unmarshal(body, &result)
		if err != nil {
			log.Fatalf("Error parsing XML: %v", err)
		}

		// Afficher les buckets de manière lisible
		if len(result.Buckets) == 0 {
			fmt.Println("No buckets found.")
			return
		}

		// Afficher les buckets
		fmt.Println("Buckets:")
		const readableDateLayout = "2006-01-02 15:04:05"
		const inputLayout = time.RFC3339 // Le format attendu de votre chaîne de date (par exemple "2024-09-17T08:58:31Z")
		
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

func init() {
	// Ajouter la commande `list-buckets` à la racine
	rootCmd.AddCommand(listBucketsCmd)
}
