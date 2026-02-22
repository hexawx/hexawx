package core

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

type PluginManager struct {
	drivers   []Driver
	exporters []Exporter
	clients   []*plugin.Client // Pour pouvoir les arrêter proprement à la fin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		drivers:   []Driver{},
		exporters: []Exporter{},
		clients:   []*plugin.Client{},
	}
}

func (m *PluginManager) Drivers() []Driver {
	return m.drivers
}

func (m *PluginManager) Exporters() []Exporter {
	return m.exporters
}

// LoadPlugin lance un binaire et l'ajoute à la liste des drivers actifs
func (m *PluginManager) LoadPlugin(path string, pluginType string) error {
	handshakeConfig := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "GOWX_PLUGIN",
		MagicCookieValue: "hello",
	}

	var p plugin.Plugin
	if pluginType == "driver" {
		p = &DriverPlugin{}
	} else {
		p = &ExporterPlugin{}
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]plugin.Plugin{
			pluginType: p,
		},
		Cmd:        exec.Command(path),
		Managed:    true,
		SyncStdout: os.Stdout, // Redirige le fmt.Print du plugin vers ta console
		SyncStderr: os.Stderr, // Redirige les erreurs vers ta console
	})

	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	raw, err := rpcClient.Dispense(pluginType)
	if err != nil {
		return err
	}

	if pluginType == "driver" {
		m.drivers = append(m.drivers, raw.(Driver))
	} else {
		m.exporters = append(m.exporters, raw.(Exporter))
	}

	m.clients = append(m.clients, client)
	return nil
}

func (m *PluginManager) StopAll() {
	for _, client := range m.clients {
		client.Kill()
	}
}
