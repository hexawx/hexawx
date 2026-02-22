package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var pluginsInstallCmd = &cobra.Command{
	Use:   "install [plugin-name]",
	Short: "Installe un plugin depuis le registry officiel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]

		// 1. Récupérer l'index
		fmt.Println("📡 Connexion au registry...")
		resp, err := http.Get("https://raw.githubusercontent.com/HexaWX/registry/main/index.json")
		if err != nil {
			fmt.Printf("❌ Impossible de contacter le registry : %v\n", err)
			return
		}
		defer resp.Body.Close()

		var catalog []RemotePlugin
		json.NewDecoder(resp.Body).Decode(&catalog)

		// 2. Chercher le plugin
		var target *RemotePlugin
		for _, p := range catalog {
			if p.Name == pluginName {
				target = &p
				break
			}
		}

		if target == nil {
			fmt.Printf("❓ Plugin '%s' introuvable dans le catalogue.\n", pluginName)
			return
		}

		// 3. Préparer le dossier de plugins
		if err := os.MkdirAll(AppConfig.Server.PluginDir, 0755); err != nil {
			fmt.Printf("❌ Impossible de créer le dossier %s : %v\n", AppConfig.Server.PluginDir, err)
			return
		}

		// 4. Préparer l'URL et le chemin
		binaryURL := resolveURL(target.BinaryURL, target.Version)
		destPath := filepath.Join("plugins", target.Name)
		if runtime.GOOS == "windows" {
			destPath += ".exe"
		}

		// 5. Téléchargement
		fmt.Printf("📥 Installation de %s (%s)...\n", target.DisplayName, target.Version)
		if err := downloadFile(destPath, binaryURL); err != nil {
			fmt.Printf("❌ Erreur de téléchargement : %v\n", err)
			return
		}

		// 6. Droits d'exécution (Linux/macOS)
		os.Chmod(destPath, 0755)
		fmt.Printf("✅ Plugin installé avec succès dans %s\n", destPath)
	},
}

// resolveURL remplace les variables {{.OS}} etc par les valeurs système
func resolveURL(template string, version string) string {
	r := strings.NewReplacer(
		"{{.Version}}", version,
		"{{.OS}}", runtime.GOOS,
		"{{.Arch}}", runtime.GOARCH,
	)
	return r.Replace(template)
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("serveur a répondu : %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func init() {
	pluginsCmd.AddCommand(pluginsInstallCmd)
}
