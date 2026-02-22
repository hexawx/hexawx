package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hexawx/hexawx/core"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Lance l'orchestrateur météo",
	Run: func(cmd *cobra.Command, args []string) {
		manager := core.NewPluginManager()

		interval := AppConfig.Server.Interval
		pluginDir := AppConfig.Server.PluginDir

		// 1. Détection et Chargement dynamique
		fmt.Println("🔍 Chargement des plugins...")

		// On scanne le dossier ./plugins
		files, _ := os.ReadDir(pluginDir)

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			path := filepath.Join("./plugins", f.Name())

			// On ne précise plus le type, on laisse le manager se débrouiller
			err := manager.AutoLoad(path, AppConfig.Plugins[f.Name()])

			if err != nil {
				fmt.Printf("❌ %s : %v\n", f.Name(), err)
			}
		}

		// 2. Nettoyage à l'arrêt
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-c
			fmt.Println("\nTerminaison propre de HexaWX...")
			manager.StopAll()
			os.Exit(0)
		}()

		// 3. Boucle de monitoring
		ticker := time.NewTicker(interval)
		fmt.Printf("📡 HexaWX est à l'écoute (Intervalle: %v) (Ctrl+C pour arrêter)\n", interval)

		for range ticker.C {
			for _, driver := range manager.Drivers() {
				data, err := driver.Fetch()
				if err != nil {
					fmt.Printf("Erreur Fetch: %v\n", err)
					continue
				}
				for _, exporter := range manager.Exporters() {
					err := exporter.Export(data)
					if err != nil {
						fmt.Printf("⚠️ Erreur Export vers un plugin : %v\n", err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
