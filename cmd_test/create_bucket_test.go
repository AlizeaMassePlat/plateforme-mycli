package cmd_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/AlizeaMassePlat/plateforme-mycli/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	// Créez un pipe pour capturer l'output
	reader, writer, _ := os.Pipe()
	defer reader.Close()

	// Sauvegarder l'ancien stdout
	oldStdout := os.Stdout
	os.Stdout = writer

	// Créez un buffer pour capturer l'output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Exécuter la fonction
	f()

	// Fermez le writer et restaurer stdout
	writer.Close()
	os.Stdout = oldStdout

	// Lire l'output capturé dans le reader
	io.Copy(&buf, reader)
	log.SetOutput(os.Stderr)

	return buf.String()
}

func TestCreateBucketCmd(t *testing.T) {
	viper.Set("s3.api_url", "http://localhost:9090")

	// Test de la création réussie du bucket
	t.Run("CreateValidBucket", func(t *testing.T) {
		output := captureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"create-bucket", "testbucket"}) 
			err := cmd.RootCmd.Execute()                                 
			assert.NoError(t, err)
		})

		// Loguer l'output pour voir ce qui est capturé
		t.Logf("Output VALID: %s", output)

		// Vérifier le message de succès
		assert.Contains(t, output, "Bucket 'testbucket' created successfully.")
	})

	// Test du cas où le bucket existe déjà
	t.Run("BucketAlreadyExists", func(t *testing.T) {
		output := captureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"create-bucket", "testbucket"}) 
			err := cmd.RootCmd.Execute()                                 
			assert.Error(t, err) // Ici on attend une erreur
		})

		// Loguer l'output pour voir ce qui est capturé
		t.Logf("Output ALREADY: %s", output)

		// Vérifier que l'output contient bien le message attendu
		assert.Contains(t, output, "Bucket already exists: Bucket already exists") 
	})
}
