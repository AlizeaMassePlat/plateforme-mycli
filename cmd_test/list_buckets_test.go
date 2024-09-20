package cmd_test

import (
	"testing"

	"github.com/AlizeaMassePlat/plateforme-mycli/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Structure pour parser la liste des buckets depuis la réponse XML
type ListAllMyBucketsResult struct {
	Buckets []struct {
		Name string `xml:"Name"`
	} `xml:"Buckets>Bucket"`
}


func TestListBucketsCmd(t *testing.T) {
	// Configurer l'URL du serveur réel
	viper.Set("s3.api_url", "http://localhost:9090") // Assurez-vous que cette URL correspond à votre serveur réel

	// Créer les buckets "list1", "list2", et "list3"
	CreateBucket(t, "list1")
	CreateBucket(t, "list2")
	CreateBucket(t, "list3")

	// Test du cas où la liste des buckets est retournée correctement
	t.Run("ListBucketsSuccess", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"list-buckets"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error during bucket listing")
		})

		// Vérifiez que les buckets créés sont bien listés
		assert.Contains(t, output, "list1", "Expected bucket 'list1' in output")
		assert.Contains(t, output, "list2", "Expected bucket 'list2' in output")
		assert.Contains(t, output, "list3", "Expected bucket 'list3' in output")
	})

	// Supprimer tous les buckets après le test de succès
	DeleteAllBuckets(t)

	// Test du cas où aucun bucket n'est trouvé
	t.Run("NoBucketsFound", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"list-buckets"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Expected no error when no buckets are found")
		})

		// Vérifiez qu'aucun bucket n'est trouvé après suppression
		assert.Contains(t, output, "No buckets found.", "Expected message when no buckets are found")
	})
}
