package core

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

func (m *PluginManager) renderHelp(term *term.Terminal) {
	format := "%-30s %s"
	helpText := fmt.Sprintf("%s", Colors.Yellow+"COMMANDES DISPONIBLES :"+Colors.Reset+"\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"list"+Colors.Reset, "Affiche les drivers et exporters chargés et leur statut.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"stop [name]"+Colors.Reset, "Arrête un plugin spécifique et libère les ressources.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"start [name]"+Colors.Reset, "Démarre un plugin spécifique.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"restart [name]"+Colors.Reset, "Arrête et redémarre un plugin spécifique.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"reload"+Colors.Reset, "Redémarre l'intégralité du système de plugins.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"install [id]"+Colors.Reset, "Installe un plugin depuis de catalogue.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"uninstall [name]"+Colors.Reset, "Déinstalle un plugin spécifique.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"stats"+Colors.Reset, "Affiche les métriques de collecte (records, erreurs).\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"catalog"+Colors.Reset, "Affiche le catalogue de plugins.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"clear"+Colors.Reset, "Efface l'écran de la console.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"help"+Colors.Reset, "Affiche ce menu d'aide.\r\n") +
		fmt.Sprintf(format, Colors.Cyan+"exit|quit"+Colors.Reset, "Ferme la session SSH.\r\n\r\n") +
		fmt.Sprintf("%s", Colors.Yellow+"Astuce: Utilisez le nom exact affiché dans 'list' pour les commandes stop/start."+Colors.Reset+"\r\n")
	term.Write([]byte(strings.ReplaceAll(helpText, "\n", "\r\n")))
}
