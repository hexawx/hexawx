package core

import "time"

type Config struct {
	Server struct {
		Interval  time.Duration `mapstructure:"interval"`
		PluginDir string        `mapstructure:"plugin_dir"`
	} `mapstructure:"server"`
	Plugins map[string]map[string]string `mapstructure:"plugins"`
}
