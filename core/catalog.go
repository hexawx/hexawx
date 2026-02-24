package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"golang.org/x/term"
)

func (m *PluginManager) getCatalog(term *term.Terminal) ([]RemotePlugin, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/hexawx/registry/main/index.json")
	if err != nil {
		term.Write([]byte(fmt.Sprintf(Prefix.Error+" Impossible de contacter le registry : %v\n", err)))
		return nil, err
	}
	defer resp.Body.Close()

	var catalog []RemotePlugin
	json.NewDecoder(resp.Body).Decode(&catalog)

	return m.filterCompatible(catalog), nil
}

func (m *PluginManager) filterCompatible(allPlugins []RemotePlugin) []RemotePlugin {
	// On récupère la plateforme locale (ex: "linux-amd64")
	currentPlatform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	compatible := make([]RemotePlugin, 0)

	for _, p := range allPlugins {
		isCompatible := false
		for _, plat := range p.SupportedPlatforms {
			if plat == currentPlatform {
				isCompatible = true
				break
			}
		}

		if isCompatible {
			compatible = append(compatible, p)
		}
	}

	return compatible
}
