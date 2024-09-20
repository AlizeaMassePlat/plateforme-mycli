package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AlizeaMassePlat/plateforme-mycli/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func TestDownloadFileCmd(t *testing.T) {
	viper.Set("s3.api_url", "http://localhost:9090")
	CreateBucket(t, "coucou2")
	CreateObject(t, "coucou2", "testdesk-1.txt", "This is the content of the test file.")
	// Créer un répertoire temporaire pour éviter les problèmes liés aux permissions
	tempDir := "."

	// Test du téléchargement réussi du fichier
	t.Run("DownloadValidFile", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"download-file", "coucou2", "testdesk-1.txt", tempDir})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error during valid file download")
		})

		// Vérifier que l'output contient le bon message
		assert.Contains(t, output, "File 'testdesk-1.txt' downloaded successfully", "Expected success message for file download")

		// Vérifier que le fichier a bien été téléchargé
		filePath := filepath.Join(tempDir, "testdesk-1.txt")
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Expected the file to be downloaded successfully")

		// Vérifier le contenu du fichier téléchargé
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, "This is the content of the test file.", string(content), "Expected file content to match")
	})

	// Test du cas où le fichier n'existe pas
	t.Run("DownloadNonExistentFile", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"download-file", "coucou", "nonexistentfile.txt", tempDir})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error when downloading a non-existent file")
		})

		// Vérifier que l'output contient le bon message d'erreur
		assert.Contains(t, output, "The system cannot find the file specified", "Expected error message for non-existent file")
	})
}
