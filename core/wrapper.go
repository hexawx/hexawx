package core

func (g *DriverRPCClient) Init(config map[string]string) error {
	var resp struct{}
	return g.client.Call("Plugin.Init", config, &resp)
}
