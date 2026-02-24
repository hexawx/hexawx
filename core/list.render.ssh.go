package core

import (
	"fmt"
	"time"

	"golang.org/x/term"
)

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
