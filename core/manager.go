package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
)

type Plugin struct {
	driver    Driver
	exporter  Exporter
	client    *plugin.Client
	config    map[string]string
	name      string
	path      string
	status    string
	version   string
	startTime time.Time
}

type PluginManager struct {
	plugins   []Plugin
	StartTime time.Time
	mu        sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{}
}

func (m *PluginManager) Plugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.plugins
}

func (m *Plugin) Driver() Driver {
	return m.driver
}

func (m *Plugin) Exporter() Exporter {
	return m.exporter
}

func (m *Plugin) Client() *plugin.Client {
	return m.client
}

func (m *Plugin) Status() string {
	return m.status
}

func (m *PluginManager) AutoLoad(config *Config) error {

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
			fmt.Printf(Prefix.Error+" %s : %v\n", f.Name(), err)
			return err
		}
	}

	return nil
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

	newPlugin := Plugin{
		client:    client,
		config:    config,
		path:      path,
		status:    Status.Running,
		version:   getVersionFromPath(path),
		startTime: time.Now(),
	}

	// 1. On tente Driver
	if raw, err := rpcClient.Dispense("driver"); err == nil && raw != nil {
		if d, ok := raw.(Driver); ok {
			d.Init(config)
			newPlugin.driver = d
			newPlugin.name, _ = d.Name()
		}
	}

	// 2. On tente Exporter seulement si Driver a échoué
	if newPlugin.driver == nil {
		if raw, err := rpcClient.Dispense("exporter"); err == nil && raw != nil {
			if e, ok := raw.(Exporter); ok {
				e.Init(config)
				newPlugin.exporter = e
				newPlugin.name, _ = e.Name()
			}
		}
	}

	if newPlugin.driver == nil && newPlugin.exporter == nil {
		client.Kill()
		return fmt.Errorf("Le plugin à l'adresse %s n'a pu être chargé ni comme driver ni comme exporter", path)
	}

	m.mu.Lock()
	m.plugins = append(m.plugins, newPlugin)
	m.mu.Unlock()

	return nil
}

func (m *PluginManager) StartPlugin(name string) (*Plugin, error) {
	var path string
	var config map[string]string
	found := false

	var currentPlugin *Plugin

	m.mu.Lock()
	for i := range m.plugins {
		if m.plugins[i].name == name {
			currentPlugin = &m.plugins[i]
			if m.plugins[i].status == Status.Running {
				return currentPlugin, fmt.Errorf("Le plugin est déjà démarré")
			}

			path = m.plugins[i].path
			config = m.plugins[i].config

			m.plugins = append(m.plugins[:i], m.plugins[i+1:]...)
			found = true
			break
		}
	}
	m.mu.Unlock()

	if !found {
		return nil, fmt.Errorf("plugin '%s' non trouvé", name)
	}

	err := m.loadPlugin(path, config)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du chargement : %w", err)
	}

	return currentPlugin, nil
}

func (m *PluginManager) StopPlugin(name string) (*Plugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.plugins {
		if m.plugins[i].name == name {
			m.plugins[i].client.Kill()
			m.plugins[i].status = Status.Stopped
			m.plugins[i].startTime = time.Time{}
			return &m.plugins[i], nil
		}
	}

	return nil, fmt.Errorf("plugin '%s' non trouvé", name)
}

func (m *PluginManager) StopAll() {
	for i := range m.plugins {
		m.StopPlugin(m.plugins[i].name)
	}
}

func (m *PluginManager) RemovePlugin(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	index := -1
	for i, p := range m.plugins {
		// On compare par le nom, c'est l'identifiant unique le plus sûr
		if p.name == name {
			index = i
			break
		}
	}

	if index != -1 {
		// Suppression et libération pour le Garbage Collector
		m.plugins = append(m.plugins[:index], m.plugins[index+1:]...)
	}
}
