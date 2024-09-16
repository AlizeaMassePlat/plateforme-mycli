package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
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
var listObjectCmd = &cobra.Command{
	Use:   "list-object",
	Short: "List objects in a specified S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Usage: list-object <bucket-name>")
		}

		bucketName := args[0]

		// Construire l'URL pour lister les objets du bucket
		url := fmt.Sprintf("http://localhost:9090/%s/", bucketName)

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
		for _, obj := range result.Objects {
			fmt.Printf("[%s] %dB %s\n", obj.LastModified, obj.Size, obj.Key)
		}
	},
}

func init() {
	rootCmd.AddCommand(listObjectCmd)
}
