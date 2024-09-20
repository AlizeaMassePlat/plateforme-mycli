package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Setup des configurations pour les tests
func setupTestConfig(t *testing.T) {
	// Charger la configuration via Viper
	viper.SetConfigName("config")  // Nom du fichier config sans l'extension
	viper.SetConfigType("yaml")    // Type du fichier de config (YAML)
	viper.AddConfigPath(".")       // Répertoire courant
	viper.AddConfigPath("../")     // Répertoire parent
	viper.AddConfigPath("/path/to/config") // Autres chemins potentiels si nécessaire

	// Lire le fichier de configuration
	err := viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Vérifier si l'URL de l'API est présente
	if viper.GetString("s3.api_url") == "" {
		t.Fatal("S3 API URL is not set in the config")
	}
}


// Fonction utilitaire pour exécuter les commandes Cobra et capturer la sortie
func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

// Test de la commande `create-bucket`
func TestCreateBucketCmd(t *testing.T) {
	setupTestConfig(t) // Charger la config depuis config.yaml

	// Nom de bucket unique pour ce test
	bucketName := fmt.Sprintf("test-bucket-%d", time.Now().Unix())

	// Exécuter la commande `create-bucket`
	_, err := executeCommand(rootCmd, "create-bucket", bucketName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Nettoyer le bucket après le test
	t.Cleanup(func() {
		_, err := executeCommand(rootCmd, "delete-bucket", bucketName)
		if err != nil {
			t.Logf("Error cleaning up bucket '%s': %v", bucketName, err)
		}
	})
}

// Test de la commande `upload-file`
func TestUploadFileCmd(t *testing.T) {
	setupTestConfig(t) // Charger la config depuis config.yaml

	// Nom du bucket
	bucketName := fmt.Sprintf("test-bucket-%d", time.Now().Unix())

	// Créer le bucket pour téléverser le fichier
	_, err := executeCommand(rootCmd, "create-bucket", bucketName)
	if err != nil {
		t.Fatalf("Failed to create bucket for upload test: %v", err)
	}

	// Chemin vers un fichier temporaire à téléverser
	tmpFile, err := os.CreateTemp("", "test-file-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Nettoyer le fichier après le test
	tmpFile.WriteString("This is a test file for upload.")
	tmpFile.Close()

	// Exécuter la commande `upload-file`
	_, err = executeCommand(rootCmd, "upload-file", bucketName, tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Nettoyer le bucket après le test
	t.Cleanup(func() {
		_, err := executeCommand(rootCmd, "delete-bucket", bucketName)
		if err != nil {
			t.Logf("Error cleaning up bucket '%s': %v", bucketName, err)
		}
	})
}

// Test de la commande `download-file`
func TestDownloadFileCmd(t *testing.T) {
	setupTestConfig(t) // Charger la config depuis config.yaml

	// Nom du bucket
	bucketName := fmt.Sprintf("test-bucket-%d", time.Now().Unix())

	// Créer le bucket pour télécharger le fichier
	_, err := executeCommand(rootCmd, "create-bucket", bucketName)
	if err != nil {
		t.Fatalf("Failed to create bucket for download test: %v", err)
	}

	// Chemin vers un fichier temporaire à téléverser
	tmpFile, err := os.CreateTemp("", "test-file-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Nettoyer le fichier après le test
	tmpFile.WriteString("This is a test file for upload.")
	tmpFile.Close()

	// Téléverser le fichier
	_, err = executeCommand(rootCmd, "upload-file", bucketName, tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to upload file for download test: %v", err)
	}

	// Télécharger le fichier
	downloadPath := fmt.Sprintf("%s-downloaded", tmpFile.Name())
	_, err = executeCommand(rootCmd, "download-file", bucketName, tmpFile.Name(), downloadPath)
	if err != nil {
		t.Fatalf("Error downloading file: %v", err)
	}

	// Vérifier que le fichier a été téléchargé avec succès
	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		t.Errorf("Expected file to be downloaded, but it was not found")
	} else {
		// Afficher un message de succès si le fichier est téléchargé
		t.Logf("File downloaded successfully to %s", downloadPath)
	}

	// Nettoyer le bucket et le fichier après le test
	t.Cleanup(func() {
		_, err := executeCommand(rootCmd, "delete-bucket", bucketName)
		if err != nil {
			t.Logf("Error cleaning up bucket '%s': %v", bucketName, err)
		}
		os.Remove(downloadPath)
	})
}


// Test de la commande `delete-bucket`
func TestDeleteBucketCmd(t *testing.T) {
	setupTestConfig(t) // Charger la config

	// Nom du bucket à supprimer
	bucketName := fmt.Sprintf("test-bucket-%d", time.Now().Unix())

	// Créer le bucket pour tester la suppression
	_, err := executeCommand(rootCmd, "create-bucket", bucketName)
	if err != nil {
		t.Fatalf("Failed to create bucket for deletion test: %v", err)
	}

	// Supprimer le bucket
	_, err = executeCommand(rootCmd, "delete-bucket", bucketName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
