package core

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// 1. Le Client RPC (Core -> Plugin)
type ExporterRPCClient struct{ client *rpc.Client }

func (g *ExporterRPCClient) Init(config map[string]string) error {
	var resp struct{}
	return g.client.Call("Plugin.Init", config, &resp)
}

func (g *ExporterRPCClient) Name() (string, error) {
	var resp string
	err := g.client.Call("Plugin.Name", struct{}{}, &resp)
	return resp, err
}

func (g *ExporterRPCClient) Export(record WeatherRecord) error {
	var resp struct{} // Obligatoire pour net/rpc même si on ne s'en sert pas
	return g.client.Call("Plugin.Export", record, &resp)
}

// 2. Le Serveur RPC (Plugin -> Core)
type ExporterRPCServer struct{ Impl Exporter }

func (s *ExporterRPCServer) Init(config map[string]string, resp *struct{}) error {
	return s.Impl.Init(config)
}

func (s *ExporterRPCServer) Name(args struct{}, resp *string) error {
	name, err := s.Impl.Name()
	if err != nil {
		return err
	}
	*resp = name
	return nil
}

func (s *ExporterRPCServer) Export(record WeatherRecord, resp *struct{}) error {
	return s.Impl.Export(record)
}

// 3. Le Plugin
type ExporterPlugin struct {
	Impl Exporter
}

func (p *ExporterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ExporterRPCServer{Impl: p.Impl}, nil
}

func (p *ExporterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ExporterRPCClient{client: c}, nil
}
