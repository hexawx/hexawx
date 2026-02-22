package core

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type ExporterRPCClient struct{ client *rpc.Client }

func (g *ExporterRPCClient) Export(record WeatherRecord) error {
	var resp struct{}
	err := g.client.Call("Plugin.Export", record, &resp)
	return err
}

type ExporterRPCServer struct{ Impl Exporter }

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
