package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

func (m *PluginManager) StartAdminShell(config *Config) {
	server := &ssh.Server{
		Addr:             fmt.Sprintf(":%d", config.Server.SshPort),
		PasswordHandler:  nil,
		PublicKeyHandler: publicKeyHandler,
		Handler: func(s ssh.Session) {

			// Affichage du MOTD
			motd := m.getColoredMOTD(s.User(), config.Server.SshPort, config.Server.Version)
			fmt.Fprint(s, motd)

			// 1. On crée un terminal interactif sur la session SSH
			term := term.NewTerminal(s, fmt.Sprintf("%s%s@%s%s$ ",
				Colors.Green,
				s.User(),
				"hexawx",
				Colors.Reset,
			))

			for {
				// 2. On lit la ligne (ReadLine gère le buffer et la touche Entrée)
				line, err := term.ReadLine()
				if err != nil {
					break
				}

				args := strings.Fields(line) // Découpe proprement (gère les espaces multiples)
				if len(args) == 0 {
					continue
				}

				cmd := args[0]

				switch cmd {
				case "list":
					m.renderList(term)
				case "stop":
					m.renderStop(term, args)
				case "start":
					m.renderStart(term, args)
				case "restart":
					m.renderRestart(term, args)
				case "reload":
					m.renderReload(term, config)
				case "install":
					m.renderInstall(term, config, args)
				case "uninstall":
					m.renderUninstall(term, args)
				case "stats":
					m.renderStats(term)
				case "catalog":
					m.renderCatalog(term)
				case "clear":
					m.renderClear(term)
				case "help":
					m.renderHelp(term)
				case "exit", "quit":
					term.Write([]byte("Fermeture de la session...\n"))
					return
				case "":
					continue
				default:
					term.Write([]byte(fmt.Sprintf("Commande inconnue : %s\n", cmd)))
					m.renderHelp(term)
				}
			}
		},
	}

	fmt.Printf("🔐 Console admin dispo sur le port %d\n", config.Server.SshPort)
	// On lance le serveur
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf(Prefix.Error+" Erreur serveur SSH : %v\n", err)
	}
}

func (m *PluginManager) renderHelp(term *term.Terminal) {
	format := "%-30s %s"
	helpText := fmt.Sprintf("%s", Colors.Yellow+"COMMANDES DISPONIBLES :"+Colors.Reset+"\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"list"+Colors.Reset, "Affiche les drivers et exporters chargés et leur statut.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"stop [name]"+Colors.Reset, "Arrête un plugin spécifique et libère les ressources.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"start [name]"+Colors.Reset, "Démarre un plugin spécifique.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"restart [name]"+Colors.Reset, "Arrête et redémarre un plugin spécifique.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"reload"+Colors.Reset, "Redémarre l'intégralité du système de plugins.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"install [id]"+Colors.Reset, "Installe un plugin depuis de catalogue.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"uninstall [name]"+Colors.Reset, "Déinstalle un plugin spécifique.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"stats"+Colors.Reset, "Affiche les métriques de collecte (records, erreurs).\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"catalog"+Colors.Reset, "Affiche le catalogue de plugins.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"clear"+Colors.Reset, "Efface l'écran de la console.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"help"+Colors.Reset, "Affiche ce menu d'aide.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"exit|quit"+Colors.Reset, "Ferme la session SSH.\r\n\r\n") +
		fmt.Sprintf("%s", Colors.Yellow+"Astuce: Utilisez le nom exact affiché dans 'list' pour les commandes stop/start."+Colors.Reset+"\r\n")
	term.Write([]byte(strings.ReplaceAll(helpText, "\n", "\r\n")))
}

func (m *PluginManager) renderClear(term *term.Terminal) {
	// \033[H  : Déplace le curseur à la position "Home" (0,0)
	// \033[2J : Efface tout l'écran
	term.Write([]byte("\033[H\033[2J"))
}

func (m *PluginManager) renderList(term *term.Terminal) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	header := fmt.Sprintf("\r\n%s%-10s %-16s %-12s %-12s %-10s %-20s%s\r\n", Colors.White, "TYPE", "NAME", "VERSION", "STATUS", "PID", "UPTIME", Colors.Reset)
	term.Write([]byte(header))
	term.Write([]byte("----------------------------------------------------------------------------------------------------\n"))

	for _, p := range m.plugins {

		pluginType := ""
		if p.driver != nil {
			pluginType = "Driver"
		} else if p.exporter != nil {
			pluginType = "Exporter"
		}

		statusColor := Colors.Red
		uptimeText := ""
		if p.status == Status.Running {
			statusColor = Colors.Green
			uptimeText = fmt.Sprintf("%s", time.Since(p.startTime).Round(time.Second))
		}

		// On formate avec des espaces fixes (%-10s) pour aligner les colonnes
		line := fmt.Sprintf("%-10s %-16s %-12s %s%-12s%s %-10s %s\r\n",
			pluginType,
			p.name,
			p.version,
			statusColor,
			p.status,
			Colors.Reset,
			p.client.ID(),
			uptimeText,
		)
		term.Write([]byte(line))

	}

	term.Write([]byte("\n\n"))
}

func (m *PluginManager) renderReload(term *term.Terminal, config *Config) {
	term.Write([]byte("Rechargement des plugins...\n"))

	m.StopAll()

	// 3. On relance l'AutoLoad
	err := m.AutoLoad(config)
	if err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" Erreur : %v\n", err)))
	} else {
		term.Write([]byte(Prefix.Success + " Plugins rechargés avec succès !\n"))
	}
}

func (m *PluginManager) renderStats(term *term.Terminal) {
	m.mu.RLock()
	count := len(m.plugins)
	m.mu.RUnlock()

	uptimeText := time.Since(m.StartTime).Round(time.Second)

	fmt.Fprintf(term, "%sMÉTRIQUES SYSTÈME :%s\r\n", Colors.Yellow, Colors.Reset)
	fmt.Fprintf(term, "  - Plugins actifs : %d\r\n", count)
	fmt.Fprintf(term, "  - Uptime         : %s\r\n", uptimeText)
}

func (m *PluginManager) renderStart(term *term.Terminal, args []string) {
	if len(args) < 2 {
		term.Write([]byte(Colors.Red + "Usage: start [Name]" + Colors.Reset + "\r\n"))
		return
	}
	target := args[1]

	// On relance le chargement (utilise une config vide ou chargée)
	if _, err := m.StartPlugin(target); err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" Lancement échoué: %v\r\n", err)))
	} else {
		term.Write([]byte(fmt.Sprintf(Prefix.Success+" Plugin %s lancé avec succès.\r\n", target)))
	}
}

func (m *PluginManager) renderStop(term *term.Terminal, args []string) {
	if len(args) < 2 {
		term.Write([]byte(Colors.Red + "Usage: stop [Name]" + Colors.Reset + "\r\n"))
		return
	}
	target := args[1]
	if _, err := m.StopPlugin(target); err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" %v"+Colors.Reset+"\r\n", err)))
	} else {
		term.Write([]byte(fmt.Sprintf(Prefix.Success+" Plugin %s arrêté.\r\n", target)))
	}
}

func (m *PluginManager) renderRestart(term *term.Terminal, args []string) {
	if len(args) < 2 {
		term.Write([]byte(Colors.Red + "Usage: restart [Name]" + Colors.Reset + "\r\n"))
		return
	}
	target := args[1]

	term.Write([]byte(fmt.Sprintf(Prefix.Restart+" Redémarrage de %s...\r\n", target)))

	_, err := m.StopPlugin(target)

	if err != nil && !strings.Contains(err.Error(), "déjà arrêté") {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" %v"+Colors.Reset+"\r\n", err)))
		return
	}
	if _, err := m.StartPlugin(target); err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" %v"+Colors.Reset+"\r\n", err)))
		return
	}
	term.Write([]byte(fmt.Sprintf(Prefix.Success+" Plugin %s relancé avec succès."+Colors.Reset+"\r\n", target)))
}

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

func (m *PluginManager) renderUninstall(term *term.Terminal, args []string) {
	if len(args) < 2 {
		term.Write([]byte(Colors.Red + "Usage: uninstall [Name]" + Colors.Reset + "\r\n"))
		return
	}
	target := args[1]

	pluginStopped, err := m.StopPlugin(target)
	if err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" %v\r\n", err)))
		return
	}

	removeFile(pluginStopped.path)
	m.RemovePlugin(target)

	term.Write([]byte(fmt.Sprintf(Prefix.Delete+" Le plugin %s a été déinstallé.\r\n", target)))
}

func (m *PluginManager) renderCatalog(term *term.Terminal) {
	catalog, err := m.getCatalog(term)
	if err != nil {
		return
	}
	header := fmt.Sprintf("\r\n%s%-10s %-20s %-20s %-10s %-10s%s\r\n", Colors.White, "TYPE", "ID", "NAME", "VERSION", "DESCRIPTION", Colors.Reset)
	term.Write([]byte(header))
	term.Write([]byte("----------------------------------------------------------------------------------------------------\n"))

	for _, c := range catalog {

		// On formate avec des espaces fixes (%-10s) pour aligner les colonnes
		line := fmt.Sprintf("%-10s %-20s %-20s %-10s %s\r\n",
			c.Type,
			c.Name,
			c.DisplayName,
			c.Version,
			c.Description,
		)
		term.Write([]byte(line))

	}

	term.Write([]byte("\n\n"))
}
