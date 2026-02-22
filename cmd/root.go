package cmd

import (
	"fmt"
	"hexawx/core"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var AppConfig core.Config

var rootCmd = &cobra.Command{
	Use:   "hexawx",
	Short: "HexaWX - Un orchestrateur météo modulaire",
	Long:  `HexaWX est un serveur de station météo basé sur une architecture hexagonale et des plugins.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de config (default: ./config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetDefault("server.interval", "5s")
	viper.SetDefault("server.plugin_dir", "./plugins")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("📖 Fichier de configuration utilisé :", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Printf("Erreur décodage config: %v\n", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
