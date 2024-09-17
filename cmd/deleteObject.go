package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		// Vérification des arguments
		if len(args) < 2 {
			log.Fatal("Usage: delete-object <bucket-name> <object-key>")
		}

		bucketName := args[0]
		objectKey := args[1]

		// Récupérer l'URL de l'API depuis la configuration via Viper
		apiURL := viper.GetString("s3.api_url")
		if apiURL == "" {
			log.Fatal("API URL is not configured. Please set it in the config file or environment variables.")
		}

		// Construire l'URL pour la requête POST de suppression
		url := fmt.Sprintf("%s/%s/?delete", apiURL, bucketName)

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
		switch resp.StatusCode {
		case http.StatusOK, http.StatusNoContent:
			fmt.Printf("Successfully deleted object '%s' from bucket '%s'.\n", objectKey, bucketName)
		case http.StatusNotFound:
			fmt.Printf("Object '%s' not found in bucket '%s'.\n", objectKey, bucketName)
		default:
			fmt.Printf("Failed to delete object '%s' from bucket '%s'. Status code: %d\n", objectKey, bucketName, resp.StatusCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteObjectCmd)
}
