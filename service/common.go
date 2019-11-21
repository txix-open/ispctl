package service

import "isp-ctl/cfg"

var (
	configClient = cfg.NewConfigClient()
	configPath   string
)

func SetConfigurationPath(path string) {
	configPath = path
}
