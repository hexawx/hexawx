package core

import (
	"fmt"

	"golang.org/x/term"
)

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
