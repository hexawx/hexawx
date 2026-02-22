package core

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type ExporterRPCClient struct{ client *rpc.Client }

func (m *ExporterRPCClient) Init(config map[string]string) error {
	var resp error
	// On appelle "Plugin.Init" sur le binaire distant
	err := m.client.Call("Plugin.Init", config, &resp)
	if err != nil {
		return err
	}
	return resp
}

func (g *ExporterRPCClient) Export(record WeatherRecord) error {
	var resp struct{}
	err := g.client.Call("Plugin.Export", record, &resp)
	return err
}

type ExporterRPCServer struct{ Impl Exporter }

func (s *ExporterRPCServer) Init(config map[string]string, resp *struct{}) error {
	return s.Impl.Init(config)
}

func (s *ExporterRPCServer) Export(record WeatherRecord, resp *struct{}) error {
	return s.Impl.Export(record)
}

type ExporterPlugin struct {
	Impl Exporter
}

func (p *ExporterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ExporterRPCServer{Impl: p.Impl}, nil
}

func (p *ExporterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ExporterRPCClient{client: c}, nil
}
