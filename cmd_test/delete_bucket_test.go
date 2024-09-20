package cmd_test

import (
	"testing"

	"github.com/AlizeaMassePlat/plateforme-mycli/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBucketCmd(t *testing.T) {
	viper.Set("s3.api_url", "http://localhost:9090")
	CreateBucket(t, "valid-bucket")
	// Test de la suppression réussie du bucket
	t.Run("DeleteValidBucket", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"delete-bucket", "valid-bucket"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error during valid bucket deletion")
		})

		assert.Contains(t, output, "Bucket 'valid-bucket' deleted successfully.", "Expected success message for bucket deletion")
	})

	// Test du cas où le bucket n'existe pas
	t.Run("DeleteNonExistentBucket", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"delete-bucket", "nonexistentbucket"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error when deleting non-existent bucket")
		})

		assert.Contains(t, output, "Bucket 'nonexistentbucket' does not exist or has already been deleted.", "Expected error message for non-existent bucket")
	})

}
