package core

import (
	"fmt"

	"golang.org/x/term"
)

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
