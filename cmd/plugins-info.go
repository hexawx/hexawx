package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pluginsListCmd = &cobra.Command{
	Use:   "info [plugin-name]",
	Short: "Affiche les information d'un plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📡 Connexion au registry...")
	},
}

func init() {
	pluginsCmd.AddCommand(pluginsListCmd)
}
