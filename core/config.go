package core

import "time"

type Config struct {
	Server struct {
		Interval  time.Duration `mapstructure:"interval"`
		PluginDir string        `mapstructure:"plugin_dir"`
		SshPort   int           `mapstructure:"ssh_port"`
	} `mapstructure:"server"`
	Plugins map[string]map[string]string `mapstructure:"plugins"`
}
