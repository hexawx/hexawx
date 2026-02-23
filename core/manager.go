package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-plugin"
)

type LoadedPlugin struct {
	Path   string
	Client *plugin.Client
}

type PluginManager struct {
	drivers       []Driver
	exporters     []Exporter
	clients       []*plugin.Client // Pour pouvoir les arrêter proprement à la fin
	activePlugins []LoadedPlugin   // Pour garder une trace des fichiers lancés
	mu            sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		drivers:   []Driver{},
		exporters: []Exporter{},
		clients:   []*plugin.Client{},
	}
}

func (m *PluginManager) Drivers() []Driver {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.drivers
}

func (m *PluginManager) Exporters() []Exporter {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.exporters
}

func (m *PluginManager) AutoLoad(config Config) error {

	pluginDir := config.Server.PluginDir

	files, _ := os.ReadDir(pluginDir)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		path := filepath.Join(pluginDir, f.Name())

		// On ne précise plus le type, on laisse le manager se débrouiller
		err := m.loadPlugin(path, config.Plugins[f.Name()])

		if err != nil {
			fmt.Printf("❌ %s : %v\n", f.Name(), err)
			return err
		}
	}

	return nil
}

func (m *PluginManager) StopAll() {
	for _, client := range m.clients {
		client.Kill()
	}
}

func (m *PluginManager) loadPlugin(path string, config map[string]string) error {
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
			m.mu.Lock()
			m.drivers = append(m.drivers, d)
			m.mu.Unlock()
			loaded = true
		}
	}

	// 2. On tente Exporter seulement si Driver a échoué
	if !loaded {
		rawExporter, err := rpcClient.Dispense("exporter")
		if err == nil && rawExporter != nil {
			if e, ok := rawExporter.(Exporter); ok {
				e.Init(config)
				m.mu.Lock()
				m.exporters = append(m.exporters, e)
				m.mu.Unlock()
				loaded = true
			}
		}
	}

	if loaded {
		m.mu.Lock()
		m.clients = append(m.clients, client)
		m.mu.Unlock()
		return nil
	}

	// Si on arrive ici, rien n'a marché
	client.Kill()
	return fmt.Errorf("Le plugin à l'adresse %s n'a pu être chargé ni comme driver ni comme exporter", path)
}
