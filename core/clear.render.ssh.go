package core

import "golang.org/x/term"

func (m *PluginManager) renderClear(term *term.Terminal) {
	// \033[H  : Déplace le curseur à la position "Home" (0,0)
	// \033[2J : Efface tout l'écran
	term.Write([]byte("\033[H\033[2J"))
}
