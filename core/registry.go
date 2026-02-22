package core

import (
	"runtime"
	"strings"
)

type RemotePlugin struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Version     string `json:"version"`
	BinaryURL   string `json:"binary_url"`
	Description string `json:"description"`
}

// GetBinaryURL remplace les placeholders par les valeurs réelles
func (p *RemotePlugin) GetBinaryURL() string {
	osName := runtime.GOOS
	archName := runtime.GOARCH

	replacer := strings.NewReplacer(
		"{{.Version}}", p.Version,
		"{{.OS}}", osName,
		"{{.Arch}}", archName,
	)
	return replacer.Replace(p.BinaryURL)
}
