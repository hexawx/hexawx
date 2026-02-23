package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pluginsInfoCmd = &cobra.Command{
	Use:   "info [plugin-name]",
	Short: "Affiche les information d'un plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📡 Connexion au registry...")
	},
}

func init() {
	pluginsCmd.AddCommand(pluginsInfoCmd)
}
