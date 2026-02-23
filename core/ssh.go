package core

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

func (m *PluginManager) StartAdminShell(config Config) {
	ssh.Handle(func(s ssh.Session) {
		// 1. On crée un terminal interactif sur la session SSH
		term := term.NewTerminal(s, "> ")

		term.Write([]byte("--- Bienvenue sur la console HexaWX ---\n"))
		term.Write([]byte("Commandes dispos : help, list, reload, exit\n\n"))

		for {
			// 2. On lit la ligne (ReadLine gère le buffer et la touche Entrée)
			line, err := term.ReadLine()
			if err != nil {
				break
			}

			cmd := strings.TrimSpace(line)

			switch cmd {
			case "help":
				term.Write([]byte("Commandes dispos : help, list, reload, exit\n\n"))
			case "list":
				m.mu.RLock()
				defer m.mu.RUnlock()

				term.Write([]byte("\r\n\033[1mTYPE       NAME            STATUS\033[0m\n"))
				term.Write([]byte("------------------------------------------\n"))

				for _, d := range m.drivers {
					name, _ := d.Name()
					// On formate avec des espaces fixes (%-10s) pour aligner les colonnes
					line := fmt.Sprintf("%-10s %-15s \033[32m[RUNNING]\033[0m\r\n", "Driver", name)
					term.Write([]byte(line))
				}

				for _, e := range m.exporters {
					name, _ := e.Name()
					line := fmt.Sprintf("%-10s %-15s \033[32m[RUNNING]\033[0m\r\n", "Exporter", name)
					term.Write([]byte(line))
				}

				term.Write([]byte("------------------------------------------\n\n"))

			case "reload":
				term.Write([]byte("Rechargement des plugins...\n"))
				m.mu.Lock()
				// 1. On arrête les anciens clients
				for _, c := range m.clients {
					c.Kill()
				}
				// 2. On vide les listes
				m.drivers = nil
				m.exporters = nil
				m.clients = nil
				m.mu.Unlock()

				// 3. On relance l'AutoLoad
				err := m.AutoLoad(config)
				if err != nil {
					term.Write([]byte(fmt.Sprintf("❌ Erreur : %v\n", err)))
				} else {
					term.Write([]byte("✅ Plugins rechargés avec succès !\n"))
				}
			case "exit", "quit":
				term.Write([]byte("Fermeture de la session...\n"))
				return
			case "":
				continue
			default:
				term.Write([]byte(fmt.Sprintf("Commande inconnue : %s\n", cmd)))
			}
		}
	})

	fmt.Printf("🔐 Console admin dispo sur le port %d\n", config.Server.SshPort)
	// On lance le serveur
	err := ssh.ListenAndServe(fmt.Sprintf(":%d", config.Server.SshPort), nil)
	if err != nil {
		fmt.Printf("❌ Erreur serveur SSH : %v\n", err)
	}
}
