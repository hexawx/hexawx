package core

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

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
