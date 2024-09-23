package cmd_test

import (
	"fmt"
	"testing"

	"github.com/AlizeaMassePlat/plateforme-mycli/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// ListBucketResult represents the response structure from the S3 API
type ListObjectResult struct {
	Objects []struct {
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		Size         int    `xml:"Size"`
	} `xml:"Objects>Object"`
}


// Test complet de la commande list-object avec un serveur réel
func TestListObjectCmd(t *testing.T) {
	// Configurer l'URL du serveur réel
	viper.Set("s3.api_url", "http://localhost:9090")

	// Nom du bucket et des objets pour le test
	bucketName := "test-list-objects"
	object1 := "test-object-1.txt"
	object2 := "test-object-2.txt"

	// Créer le bucket
	CreateBucket(t, bucketName)

	// Ajouter des objets au bucket
	CreateObject(t, bucketName, object1, "This is the content of test object 1.")
	CreateObject(t, bucketName, object2, "This is the content of test object 2.")

	// Lister les objets et capturer la sortie
	t.Run("ListObjectsInBucket", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"list-object", bucketName})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error during bucket listing")
		})
		// Vérifier que les objets sont listés
		assert.Contains(t, output, object1, fmt.Sprintf("Expected object %s to be listed", object1))
		assert.Contains(t, output, object2, fmt.Sprintf("Expected object %s to be listed", object2))
	})

	// Supprimer le bucket et ses objets
	DeleteBucket(t, bucketName)
	CreateBucket(t, "test-sans-objet")

	t.Run("NoObjectsFound", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"list-object", "test-sans-objet"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error when no object are found")
		})

		// Vérifier qu'aucun bucket n'est trouvé après suppression
		assert.Contains(t, output, "No objects found.", "Expected message when no buckets are found")
	})

}
