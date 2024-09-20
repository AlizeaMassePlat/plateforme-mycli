package cmd_test

import (
    "os"
    "testing"
    "github.com/AlizeaMassePlat/plateforme-mycli/cmd"
    "github.com/spf13/viper"
    "github.com/stretchr/testify/assert"
)

func setupTestFile(t *testing.T, filePath string, content string) {
    err := os.WriteFile(filePath, []byte(content), 0644)
    if err != nil {
        t.Fatalf("Failed to create test file: %v", err)
    }
}

func TestUploadFileCmd(t *testing.T) {
    viper.Set("s3.api_url", "http://localhost:9090")
    CreateBucket(t, "upload-bucket")
    file := "./testdesk-1.txt"
    // Créez le fichier de test avec un contenu spécifique
    setupTestFile(t, file, "This is the content of the test file.")
    
    t.Run("UploadValidFile", func(t *testing.T) {
        defer os.Remove(file) // Nettoyage après le test

        output := CaptureOutput(func() {
            cmd.RootCmd.SetArgs([]string{"upload-file", "upload-bucket", file})
            err := cmd.RootCmd.Execute()
            assert.NoError(t, err, "Expected no error during valid file upload")
        })

        t.Logf("Captured output for valid file upload: %s", output)

        // Nettoyer l'output en supprimant les caractères indésirables

        // Vérifiez que l'output contient le bon message
        assert.Contains(t, output, "File 'testdesk-1.txt' uploaded successfully to bucket 'upload-bucket'.")
    })

}





