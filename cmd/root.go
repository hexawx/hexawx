package cmd

import (
	"fmt"
	"os"

	"github.com/hexawx/hexawx/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var AppConfig *core.Config

var rootCmd = &cobra.Command{
	Use:   "hexawx",
	Short: "HexaWX - Station météo modulaire",
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
	viper.SetDefault("server.ssh_port", 2233)

	viper.ReadInConfig()

	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Printf("Erreur décodage config: %v\n", err)
	}

	AppConfig.Server.Version = "1.0.0"
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
