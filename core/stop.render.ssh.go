package core

import (
	"fmt"

	"golang.org/x/term"
)

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
