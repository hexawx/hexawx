package cmd

import (
	"fmt"
	"os"
	"os/signal"
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
		manager.StartTime = time.Now()

		interval := AppConfig.Server.Interval

		go func() {
			manager.StartAdminShell(AppConfig)
		}()

		// 1. Détection et Chargement dynamique
		fmt.Println("🔍 Chargement des plugins...")

		// On scanne le dossier ./plugins
		manager.AutoLoad(AppConfig)

		// 2. Nettoyage à l'arrêt
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-c
			fmt.Println("\nFermeture propre de HexaWX...")
			manager.StopAll()
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}()

		// 3. Boucle de monitoring
		ticker := time.NewTicker(interval)
		fmt.Printf("\n📡 HexaWX est à l'écoute (Intervalle: %v) (Ctrl+C pour arrêter)\n\n", interval)

		for range ticker.C {
			for _, pluginDriver := range manager.Plugins() {
				if pluginDriver.Driver() == nil || pluginDriver.Status() != core.Status.Running {
					continue
				}

				data, err := pluginDriver.Driver().Fetch()
				if err != nil {
					fmt.Printf("Erreur Fetch: %v\n", err)
					continue
				}

				for _, pluginExporter := range manager.Plugins() {
					if pluginExporter.Exporter() == nil || pluginExporter.Status() != core.Status.Running {
						continue
					}

					err := pluginExporter.Exporter().Export(data)
					if err != nil {
						fmt.Printf(core.Prefix.Warning+" Erreur Export vers un plugin : %v\n", err)
					}

				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
