package cmd_test

import (
	"testing"
	"github.com/AlizeaMassePlat/plateforme-mycli/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func TestDeleteObjectCmd(t *testing.T) {

	viper.Set("s3.api_url", "http://localhost:9090")
	CreateBucket(t, "coucou")
	CreateObject(t, "coucou", "valid-object.txt", "coucoutestcontent")
	// Test de la suppression réussie de l'objet
	t.Run("DeleteValidObject", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"delete-object", "coucou", "valid-object.txt"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Unexpected error during valid object deletion")
		})

		assert.Contains(t, output, "Successfully deleted object 'valid-object.txt' from bucket 'coucou'.", "Expected success message for object deletion")
	})

	// Test du cas où l'objet n'existe pas
	t.Run("DeleteNonExistentObject", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"delete-object", "coucou", "nonexistent-object"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Unexpected error when deleting non-existent object")
		})

		assert.Contains(t, output, "Object 'nonexistent-object' not found in bucket 'coucou'.", "Expected error message for non-existent object")
	})

}
