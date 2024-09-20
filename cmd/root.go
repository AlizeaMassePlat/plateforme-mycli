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
var RootCmd = &cobra.Command{
	Use:   "bs3",
	Short: "CLI to interact with S3 API",
}

// Execute exécute la commande root et toutes ses sous-commandes
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init initie la configuration
func init() {
	cobra.OnInitialize(initConfig)

	// Définir un flag pour permettre à l'utilisateur de spécifier un fichier de configuration
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-cli.yaml)")
}

// initConfig configure Viper pour lire les fichiers de configuration et les variables d'environnement
func initConfig() {
	// Vérifier si la variable d'environnement MYCLI_CONFIG est définie
	envConfig := os.Getenv("MYCLI_CONFIG")

	if cfgFile != "" {
		// Utiliser le fichier spécifié par l'utilisateur avec --config
		viper.SetConfigFile(cfgFile)
	} else if envConfig != "" {
		// Utiliser le fichier spécifié dans la variable d'environnement MYCLI_CONFIG
		viper.SetConfigFile(envConfig)
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
	} else if cfgFile != "" || envConfig != "" {
		log.Fatalf("Could not read config file: %s", viper.ConfigFileUsed())
	}
}
