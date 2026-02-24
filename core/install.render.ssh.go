package core

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func (m *PluginManager) renderInstall(term *term.Terminal, config *Config, args []string) {
	if len(args) < 2 {
		term.Write([]byte(Colors.Red + "Usage: install [Id]" + Colors.Reset + "\r\n"))
		return
	}
	pluginName := args[1]

	// 1. Récupérer l'index
	catalog, err := m.getCatalog(term)
	if err != nil {
		return
	}

	// 2. Chercher le plugin
	var target *RemotePlugin
	for _, p := range catalog {
		if p.Name == pluginName {
			target = &p
			break
		}
	}

	if target == nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" Plugin '%s' introuvable dans le catalogue.\r\n", pluginName)))
		return
	}

	// 3. Préparer le dossier de plugins
	if err := os.MkdirAll(config.Server.PluginDir, 0755); err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" Impossible de créer le dossier %s : %v\r\n", config.Server.PluginDir, err)))
		return
	}

	// 4. Préparer l'URL et le chemin
	binaryURL, destPath := resolveURL(target.BinaryURL, target.Version, config.Server.PluginDir, target.Name)

	// 5. Téléchargement
	term.Write([]byte(fmt.Sprintf(Prefix.Install+" Plugin %s (%s)...\r\n", target.DisplayName, target.Version)))
	if err := downloadFile(destPath, binaryURL); err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" Erreur de téléchargement : %v\r\n", err)))
		return
	}

	// 6. Droits d'exécution (Linux/macOS)
	os.Chmod(destPath, 0755)

	// 7. Démarrage du plugin
	if err := m.loadPlugin(destPath, nil); err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" Erreur d'activation : %v\r\n", err)))
		return
	}

	term.Write([]byte(fmt.Sprintf(Prefix.Success+" L'installation du plugin %s est terminée.\r\n", target.DisplayName)))
}
