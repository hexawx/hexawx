package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var pluginsInfoCmd = &cobra.Command{
	Use:   "list",
	Short: "Liste les plugins depuis le registry officiel",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔎 Recherche des plugins disponibles...")

		resp, err := http.Get("https://raw.githubusercontent.com/HexaWX/registry/main/index.json")
		if err != nil {
			fmt.Println("❌ Erreur : Impossible de joindre le catalogue.")
			return
		}
		defer resp.Body.Close()

		var catalog []RemotePlugin
		json.NewDecoder(resp.Body).Decode(&catalog)

		fmt.Printf("\n%-20s %-10s %-30s\n", "NOM", "VERSION", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 65))

		for _, p := range catalog {
			fmt.Printf("%-20s %-10s %-30s\n", p.Name, p.Version, p.Description)
		}
	},
}

func init() {
	pluginsCmd.AddCommand(pluginsInfoCmd)
}
