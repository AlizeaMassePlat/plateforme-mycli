package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type DeleteObjectRequest struct {
	XMLName xml.Name `xml:"Delete"`
	Objects []ObjectToDelete `xml:"Object"`
}

type ObjectToDelete struct {
	Key string `xml:"Key"`
}

// deleteObjectCmd represents the deleteObject command
var deleteObjectCmd = &cobra.Command{
	Use:   "delete-object",
	Short: "Deletes an object from the specified S3 bucket",
	Long: `This command deletes an object from the specified S3 bucket.
You need to specify the bucket name and the object key.
For example:

my-cli delete-object <bucket-name> <object-key>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Usage: delete-object <bucket-name> <object-key>")
			return
		}

		bucketName := args[0]
		objectKey := args[1]

		// Construire l'URL pour l'API
		url := fmt.Sprintf("http://localhost:9090/%s/?delete", bucketName)

		// Préparer la requête de suppression en XML
		deleteReq := DeleteObjectRequest{
			Objects: []ObjectToDelete{
				{Key: objectKey},
			},
		}

		// Encoder la requête en XML
		xmlData, err := xml.Marshal(deleteReq)
		if err != nil {
			log.Fatalf("Error marshalling XML: %v", err)
		}

		// Créer la requête HTTP POST
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlData))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		// Ajouter les en-têtes nécessaires
		req.Header.Set("Content-Type", "application/xml")

		// Envoyer la requête
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		// Traiter la réponse
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Successfully deleted object '%s' from bucket '%s'.\n", objectKey, bucketName)
		} else {
			fmt.Printf("Failed to delete object. Status code: %d\n", resp.StatusCode)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteObjectCmd)
}
