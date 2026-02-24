package core

import (
	"fmt"
	"strconv"
	"strings"
)

var MOTD = BoldColors.Cyan + `┌─────────────────────────────────────────────────────┐` + "\r\n" +
	BoldColors.Cyan + `│  _    _                      __          __ __   __ │` + "\r\n" +
	BoldColors.Cyan + `│ | |  | |                     \ \        / / \ \ / / │` + "\r\n" +
	BoldColors.Cyan + `│ | |__| |  ___ __  __  __ _    \ \  /\  / /   \   /  │` + "\r\n" +
	BoldColors.Cyan + `│ |  __  | / _ \\ \/ / / _` + "`" + ` |    \ \/  \/ /     > <   │` + "\r\n" +
	BoldColors.Cyan + `│ | |  | ||  __/ >  < | (_| |     \  /\  /     / ^ \  │` + "\r\n" +
	BoldColors.Cyan + `│ |_|  |_| \___|/_/\_\ \__,_|      \/  \/     /_/ \_\ │` + "\r\n" +
	BoldColors.Cyan + `│                                                     │` + "\r\n" +
	BoldColors.Cyan + `│             ` + Colors.Yellow + `>> HEXAWX SYSTEM CONSOLE <<` + BoldColors.Cyan + `             │` + "\r\n" +
	BoldColors.Cyan + `│                                                     │` + "\r\n" +
	BoldColors.Cyan + `│ ` + Colors.Reset + ` Version: %-17s` + Colors.Green + `Status: Operational` + BoldColors.Cyan + `      │` + "\r\n" +
	BoldColors.Cyan + `│ ` + Colors.Reset + ` Admin Port: %-14s` + Colors.Reset + `User: %-19s` + BoldColors.Cyan + `│` + "\r\n" +
	BoldColors.Cyan + `└─────────────────────────────────────────────────────┘` + Colors.Reset + "\r\n"

func (m *PluginManager) getColoredMOTD(user string, port int, version string) string {
	sPort := strconv.Itoa(port)
	return strings.ReplaceAll(fmt.Sprintf(MOTD, version, sPort, user), "\n", "\r\n")
}
