package core

import (
	"fmt"

	"golang.org/x/term"
)

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
