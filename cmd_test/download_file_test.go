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
	dir := "."

	// Test du téléchargement réussi du fichier
	t.Run("DownloadValidFile", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"download-file", "coucou2", "testdesk-1.txt", dir})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error during valid file download")
		})

		// Vérifier que l'output contient le bon message
        assert.Contains(t, output, "Download completed successfully.")

		// Vérifier que le fichier a bien été téléchargé
		filePath := filepath.Join(dir, "testdesk-1.txt")
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
			cmd.RootCmd.SetArgs([]string{"download-file", "coucou", "nonexistentfile.txt", dir})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error when downloading a non-existent file")
		})

		// Vérifier que l'output contient le bon message d'erreur
		assert.Contains(t, output, "the system cannot find the file specified")

		
	})
}
