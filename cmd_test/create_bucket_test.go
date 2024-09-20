package cmd_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"log"
	"fmt"
    "net/http"
	"github.com/AlizeaMassePlat/plateforme-mycli/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"encoding/xml"
	"strings"
)

func CaptureOutput(f func()) string {
    // Création du pipe pour capturer la sortie
    reader, writer, err := os.Pipe()
    if err != nil {
        log.Fatalf("failed to create pipe: %v", err) // Utiliser log.Fatal pour arrêter l'exécution en cas d'erreur
    }
    defer reader.Close() // Assurez-vous de toujours fermer le reader

    // Sauvegarder stdout et le writer de logs
    oldStdout := os.Stdout
    oldLog := log.Writer()

    // Rediriger stdout et les logs vers le writer
    os.Stdout = writer
    log.SetOutput(writer)

    // Buffer pour capturer la sortie
    var buf bytes.Buffer

    // Exécution de la fonction
    done := make(chan error, 1)
    go func() {
        _, err := io.Copy(&buf, reader) // Copier la sortie dans le buffer
        done <- err
    }()

    f()

    // Restaurer stdout et les logs
    writer.Close()
    os.Stdout = oldStdout
    log.SetOutput(oldLog)

    // Attendre la copie des données
    if err := <-done; err != nil {
        log.Fatalf("failed to copy output: %v", err) // Utiliser log.Fatal pour arrêter l'exécution en cas d'erreur
    }

    return buf.String()
}


func CreateBucket(t *testing.T, bucketName string) {
	apiURL := "http://localhost:9090" // URL de votre serveur S3-like
	createURL := fmt.Sprintf("%s/%s/", apiURL, bucketName)

	// Créer la requête PUT pour créer le bucket
	req, err := http.NewRequest("PUT", createURL, nil)
	assert.NoError(t, err, "Failed to create PUT request for bucket creation")

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Failed to create bucket")
	defer resp.Body.Close()

	// Vérifier que le bucket a bien été créé
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("Expected status 200 when creating bucket %s", bucketName))
}

func DeleteBucket(t *testing.T, bucketName string) {
	apiURL := "http://localhost:9090" // URL de votre serveur S3-like
	deleteURL := fmt.Sprintf("%s/%s/", apiURL, bucketName)

	// Créer la requête DELETE pour supprimer le bucket
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	assert.NoError(t, err, "Failed to create DELETE request for bucket deletion")

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Failed to delete bucket")
	defer resp.Body.Close()

	// Si le bucket n'existe pas, ignorer l'erreur 404
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("Failed to delete bucket %s: status code %d", bucketName, resp.StatusCode)
	}
}

func DeleteAllBuckets(t *testing.T) {
	apiURL := "http://localhost:9090" // URL de votre serveur S3-like

	// Envoyer une requête GET pour lister les buckets
	resp, err := http.Get(apiURL)
	assert.NoError(t, err, "Failed to list buckets")
	defer resp.Body.Close()

	// Vérifier le code de réponse
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 when listing buckets")

	// Lire et parser la réponse XML
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")

	var result ListAllMyBucketsResult
	err = xml.Unmarshal(body, &result)
	assert.NoError(t, err, "Failed to parse XML response")

	// Supprimer chaque bucket
	for _, bucket := range result.Buckets {
		DeleteBucket(t, bucket.Name) 
	}
}

func CreateObject(t *testing.T, bucketName, objectName, objectContent string) {
	apiURL := "http://localhost:9090" // URL de votre serveur S3-like
	objectURL := fmt.Sprintf("%s/%s/%s", apiURL, bucketName, objectName)

	// Créer une requête PUT pour ajouter l'objet dans le bucket
	req, err := http.NewRequest("PUT", objectURL, strings.NewReader(objectContent))
	assert.NoError(t, err, "Failed to create PUT request for object creation")

	// Ajouter l'en-tête Content-Type (par exemple, text/plain ou autre selon le type de contenu)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("X-Amz-Decoded-Content-Length", "12")
	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Failed to create object in bucket")
	defer resp.Body.Close()

	// Vérifier que l'objet a bien été créé (status code 200)
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("Expected status 200 when creating object %s in bucket %s", objectName, bucketName))
}


func TestCreateBucketCmd(t *testing.T) {
	viper.Set("s3.api_url", "http://localhost:9090")

	// Test de la création réussie du bucket
	t.Run("CreateValidBucket", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"create-bucket", "testbucket"})
			err := cmd.RootCmd.Execute() 
			assert.NoError(t, err, "Unexpected error during valid bucket creation")
		})
		
	
		assert.Contains(t, output, "Bucket 'testbucket' created successfully.", "Expected success message for bucket creation")
	})
	

	// Test du cas où le bucket existe déjà
	t.Run("BucketAlreadyExists", func(t *testing.T) {
		output := CaptureOutput(func() {
			cmd.RootCmd.SetArgs([]string{"create-bucket", "testbucket"})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err, "Unexpected error when bucket already exists")
		})
		

		assert.Contains(t, output, "Bucket 'testbucket' already exists", "Expected error message for existing bucket")
	})
}
