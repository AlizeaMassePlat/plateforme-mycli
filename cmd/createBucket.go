package cmd

import (
    "fmt"
    "net/http"
    "log"
    "github.com/spf13/cobra"
)

// createBucketCmd représente la commande create-bucket
var createBucketCmd = &cobra.Command{
    Use:   "create-bucket",
    Short: "Create a new S3 bucket via the API",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Fatal("Bucket name is required")
        }
        bucketName := args[0]
        
        // URL de votre API qui permet de créer un bucket
        url := "http://localhost:9090/" + bucketName + "/"

        // Faire une requête PUT pour créer un bucket
        req, err := http.NewRequest("PUT", url, nil)
        if err != nil {
            log.Fatalf("Error creating request: %v", err)
        }

        // Envoyer la requête
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Fatalf("Error making request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            fmt.Printf("Bucket '%s' created successfully.\n", bucketName)
        } else {
            fmt.Printf("Failed to create bucket. Status code: %d\n", resp.StatusCode)
        }
    },
}

func init() {
    rootCmd.AddCommand(createBucketCmd)
}
