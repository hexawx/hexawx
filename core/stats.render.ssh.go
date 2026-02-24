package core

import (
	"fmt"
	"runtime"
	"time"

	"golang.org/x/term"
)

func (m *PluginManager) renderStats(term *term.Terminal) {
	m.mu.RLock()
	count := len(m.plugins)
	m.mu.RUnlock()

	uptimeText := time.Since(m.StartTime).Round(time.Second)

	fmt.Fprintf(term, "%sCARACTÉRISTIQUES SYSTÈME :%s\r\n", Colors.Yellow, Colors.Reset)
	fmt.Fprintf(term, "  - OS : %s\r\n", runtime.GOOS)
	fmt.Fprintf(term, "  - Arch : %s\r\n", runtime.GOARCH)

	fmt.Fprintf(term, "\r\n%sMÉTRIQUES SYSTÈME :%s\r\n", Colors.Yellow, Colors.Reset)
	fmt.Fprintf(term, "  - Plugins actifs : %d\r\n", count)
	fmt.Fprintf(term, "  - Uptime         : %s\r\n", uptimeText)
}
