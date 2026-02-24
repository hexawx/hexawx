package core

import (
	"fmt"
	"strings"

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
