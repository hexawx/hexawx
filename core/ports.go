package core

import (
	"time"

	"github.com/hashicorp/go-plugin"
)

// La donnée météo universelle pour HexaWX
type WeatherRecord struct {
	Timestamp   time.Time
	Temperature float64
	Humidity    float64
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "HEXAWX_PLUGIN",
	MagicCookieValue: "hello",
}

// L'interface qu'un plugin "Driver" doit implémenter
// Driver est le port d'entrée (Infrastructure -> Core)
type Driver interface {
	Init(config map[string]string) error
	Fetch() (WeatherRecord, error)
}

// L'interface qu'un plugin "Exporter" doit implémenter
// Exporter est le port de sortie (Core -> Infrastructure)
type Exporter interface {
	Init(config map[string]string) error
	Export(record WeatherRecord) error
}
