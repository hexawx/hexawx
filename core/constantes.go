package core

var Status = struct {
	Running string
	Stopped string
}{
	Running: "RUNNING",
	Stopped: "STOPPED",
}

var Prefix = struct {
	Success     string
	Error       string
	Warning     string
	Information string
	Restart     string
	Install     string
	Delete      string
}{
	Success:     Colors.Green + "[✅ Succès]" + Colors.Reset,
	Error:       Colors.Red + "[❌ Erreur]" + Colors.Reset,
	Warning:     Colors.Yellow + "[⚠️ Attention]" + Colors.Reset,
	Information: Colors.Blue + "[ℹ️ Info]" + Colors.Reset,
	Restart:     Colors.Cyan + "[🔄 Redémarrage]" + Colors.Reset,
	Install:     Colors.Green + "[📥 Installation]" + Colors.Reset,
	Delete:      Colors.Red + "[🗑️ Suppression]" + Colors.Reset,
}

var Colors = struct {
	Blue   string
	Cyan   string
	Green  string
	Red    string
	Reset  string
	White  string
	Yellow string
}{
	Blue:   "\033[34m",
	Cyan:   "\033[36m",
	Green:  "\033[32m",
	Red:    "\033[31m",
	Reset:  "\033[0m",
	White:  "\033[1m",
	Yellow: "\033[33m",
}

var BoldColors = struct {
	Blue   string
	Cyan   string
	Green  string
	Red    string
	Reset  string
	White  string
	Yellow string
}{
	Blue:   "\033[1;34m",
	Cyan:   "\033[1;36m",
	Green:  "\033[1;32m",
	Red:    "\033[1;31m",
	Reset:  "\033[1;0m",
	White:  "\033[1;1m",
	Yellow: "\033[1;33m",
}
