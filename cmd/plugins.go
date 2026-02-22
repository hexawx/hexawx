package cmd

import (
	"github.com/spf13/cobra"
)

// RemotePlugin correspond à la structure de ton index.json
type RemotePlugin struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Version     string `json:"version"`
	BinaryURL   string `json:"binary_url"`
	Description string `json:"description"`
}

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Gestion des plugins (install, list, info)",
	Long:  `Permet de gérer l'écosystème de plugins HexaWX : recherche, installation et mise à jour.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Si l'utilisateur tape juste "hexawx plugins", on affiche l'aide
		cmd.Help()
	},
}

func init() {
	// On attache "plugins" à la commande racine de ton application (rootCmd)
	rootCmd.AddCommand(pluginsCmd)
}
