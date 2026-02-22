package core

import (
	"fmt"
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

func (m *PluginManager) AutoLoad(path string, config map[string]string) error {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"driver":   &DriverPlugin{},
			"exporter": &ExporterPlugin{},
		},
		Cmd: exec.Command(path),
	})

	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// 1. On tente de le récupérer comme un Driver
	rawDriver, err := rpcClient.Dispense("driver")
	if err == nil && rawDriver != nil {
		d := rawDriver.(Driver)
		d.Init(config)
		m.drivers = append(m.drivers, d)
		m.clients = append(m.clients, client)
		fmt.Printf("🔌 Driver détecté et chargé : %s\n", path)
		return nil
	}

	// 2. Sinon, on tente comme un Exporter
	rawExporter, err := rpcClient.Dispense("exporter")
	if err == nil && rawExporter != nil {
		e := rawExporter.(Exporter)
		e.Init(config)
		m.exporters = append(m.exporters, e)
		m.clients = append(m.clients, client)
		fmt.Printf("📦 Exporter détecté et chargé : %s\n", path)
		return nil
	}

	client.Kill()
	return fmt.Errorf("type de plugin inconnu")
}

// LoadPlugin lance un binaire et l'ajoute à la liste des drivers actifs
func (m *PluginManager) LoadPlugin(path string, pluginType string) error {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"driver":   &DriverPlugin{},
			"exporter": &ExporterPlugin{},
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
