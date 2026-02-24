package core

import (
	"fmt"

	"golang.org/x/term"
)

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
