package cmd

import (
	"fmt"
	"log"
	"net/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteBucketCmd représente la commande `delete-bucket`
var DeleteBucketCmd = &cobra.Command{
	Use:   "delete-bucket",
	Short: "Delete an S3 bucket via the API",
	Run: func(cmd *cobra.Command, args []string) {
		// Vérification que le nom du bucket est fourni
		if len(args) < 1 {
			log.Fatal("Bucket name is required")
		}
		bucketName := args[0]

		// Récupérer l'URL du serveur S3 à partir de la configuration
		apiURL := viper.GetString("s3.api_url")
		if apiURL == "" {
			log.Fatal("API URL is not configured. Please set it in the config file or environment variables.")
		}

		// Construire l'URL pour la requête DELETE
		url := fmt.Sprintf("%s/%s/", apiURL, bucketName)

		// Créer la requête HTTP DELETE
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatalf("Error creating DELETE request: %v", err)
		}

		// Envoyer la requête DELETE
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making DELETE request: %v", err)
		}
		defer resp.Body.Close()

		// Vérifier le statut de la réponse
		switch resp.StatusCode {
		case http.StatusOK, http.StatusNoContent:
			fmt.Printf("Bucket '%s' deleted successfully.\n", bucketName)
		case http.StatusNotFound:
			fmt.Printf("Bucket '%s' does not exist or has already been deleted.\n", bucketName)
		case http.StatusInternalServerError: 
			fmt.Printf("Internal server error : Status code: %d\n", resp.StatusCode)
		default:
			fmt.Printf("Failed to delete bucket '%s'. Status code: %d\n", bucketName, resp.StatusCode)
		}
	},
}

func init() {
	// Enregistrer la commande delete-bucket dans la racine
	RootCmd.AddCommand(DeleteBucketCmd)
}
