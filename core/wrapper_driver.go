package core

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// 1. Le Client RPC : Comment le Serveur appelle le Plugin
type DriverRPCClient struct{ client *rpc.Client }

func (g *DriverRPCClient) Fetch() (WeatherRecord, error) {
	var resp WeatherRecord
	err := g.client.Call("Plugin.Fetch", new(interface{}), &resp)
	return resp, err
}

// 2. Le Serveur RPC : Comment le Plugin répond au Serveur
type DriverRPCServer struct{ Impl Driver }

func (s *DriverRPCServer) Fetch(args interface{}, resp *WeatherRecord) error {
	data, err := s.Impl.Fetch()
	*resp = data
	return err
}

// 3. L'implémentation de go-plugin
type DriverPlugin struct {
	Impl Driver
}

func (p *DriverPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DriverRPCServer{Impl: p.Impl}, nil
}

func (p *DriverPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DriverRPCClient{client: c}, nil
}
