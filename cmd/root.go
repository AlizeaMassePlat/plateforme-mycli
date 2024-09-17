package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd représente la commande de base
var rootCmd = &cobra.Command{
	Use:   "my-cli",
	Short: "CLI to interact with S3 API",
}

// Execute exécute la commande root et toutes ses sous-commandes
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init initie la configuration
func init() {
	cobra.OnInitialize(initConfig)

	// Définir un flag pour permettre à l'utilisateur de spécifier un fichier de configuration
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-cli.yaml)")
}

// initConfig configure Viper pour lire les fichiers de configuration et les variables d'environnement
func initConfig() {
	if cfgFile != "" {
		// Utiliser le fichier spécifié par l'utilisateur
		viper.SetConfigFile(cfgFile)
	} else {
		// Utiliser le fichier de configuration par défaut dans le répertoire HOME
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".my-cli")
	}

	// Lire les variables d'environnement qui commencent par MYCLI_
	viper.SetEnvPrefix("mycli")
	viper.AutomaticEnv()

	// Lire le fichier de configuration
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if cfgFile != "" {
		log.Fatalf("Could not read config file: %s", cfgFile)
	}
}
