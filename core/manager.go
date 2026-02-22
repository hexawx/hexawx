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
		Cmd:        exec.Command(path),
		Managed:    true,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,
	})

	rpcClient, err := client.Client()
	if err != nil {
		return fmt.Errorf("erreur client RPC: %w", err)
	}

	// On utilise une variable pour suivre si on a réussi à charger quelque chose
	var loaded bool

	// 1. On tente Driver
	rawDriver, err := rpcClient.Dispense("driver")
	if err == nil && rawDriver != nil {
		if d, ok := rawDriver.(Driver); ok {
			d.Init(config)
			m.drivers = append(m.drivers, d)
			loaded = true
		}
	}

	// 2. On tente Exporter seulement si Driver a échoué
	if !loaded {
		rawExporter, err := rpcClient.Dispense("exporter")
		if err == nil && rawExporter != nil {
			if e, ok := rawExporter.(Exporter); ok {
				e.Init(config)
				m.exporters = append(m.exporters, e)
				loaded = true
			}
		}
	}

	if loaded {
		m.clients = append(m.clients, client)
		return nil
	}

	// Si on arrive ici, rien n'a marché
	client.Kill()
	return fmt.Errorf("Le plugin à l'adresse %s n'a pu être chargé ni comme driver ni comme exporter", path)
}

func (m *PluginManager) StopAll() {
	for _, client := range m.clients {
		client.Kill()
	}
}
