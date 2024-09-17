package cmd

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "io"
    "github.com/spf13/cobra"
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

        // URL de votre API qui permet de créer un bucket
        url := "http://localhost:9090/" + bucketName + "/"

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
    if resp.StatusCode == http.StatusOK {
        return nil
    } else if resp.StatusCode == http.StatusConflict {
        return fmt.Errorf("%s", string(body))
    } else {
        return fmt.Errorf("unexpected error: received status code %d with message: %s", resp.StatusCode, string(body))
    }
}

func init() {
    rootCmd.AddCommand(createBucketCmd)
}
