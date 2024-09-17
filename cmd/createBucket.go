package cmd

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "io"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// createBucketCmd représente la commande create-bucket
var createBucketCmd = &cobra.Command{
    Use:   "create-bucket",
    Short: "Create a new S3 bucket via the API",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Println("Error: Bucket name is required")
            return
        }

        bucketName := args[0]

        // Utilisation de Viper pour récupérer l'URL de l'API depuis une variable d'environnement ou un fichier de config
        apiURL := viper.GetString("s3.api_url")
        if apiURL == "" {
            log.Println("Error: S3 API URL is not configured. Please set it in the config file or environment variables.")
            return
        }

        // URL complète pour créer le bucket
        url := fmt.Sprintf("%s/%s/", apiURL, bucketName)

        // Appel pour créer le bucket
        if err := createBucket(url); err != nil {
            log.Printf("Error creating bucket: %v\n", err)
        } else {
            fmt.Printf("Bucket '%s' created successfully.\n", bucketName)
        }
    },
}

// Fonction de création de bucket avec gestion des erreurs
func createBucket(url string) error {
    // Créer une requête PUT pour créer le bucket
    req, err := http.NewRequest("PUT", url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    // Créer un client HTTP avec un timeout pour éviter les requêtes bloquantes
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Envoyer la requête
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    // Lire le corps de la réponse
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }

    // Vérifier le code de statut HTTP et afficher un message détaillé
    switch resp.StatusCode {
        case http.StatusOK: // 200 - Succes
            return nil
        case http.StatusConflict: // 409 - Bucket already exists
            return fmt.Errorf("%s", string(body))
        case http.StatusBadRequest: // 400 - Invalid bucket name
            return fmt.Errorf("invalid bucket name or bad request: %s", string(body))
        case http.StatusNotFound: // 404 - API endpoint not found
            return fmt.Errorf("API endpoint not found: %s", string(body))
        case http.StatusInternalServerError: // 500 - Server error
            return fmt.Errorf("server error: %s", string(body))
        default:
            return fmt.Errorf("unexpected error: received status code %d with message: %s", resp.StatusCode, string(body))
    }
}

func init() {
    rootCmd.AddCommand(createBucketCmd)
}
