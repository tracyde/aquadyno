package config

import (
	"github.com/tracyde/aquadyno/probe"
)

type Config struct {
	ConfigManager string
	ProbeManager  string
	UserName      string
	ApiKey        string
	UpdateUrl     string
	UpdateMethod  string
	Probes        []probe.Probe
}

func NewDefaultConfig() *Config {
	return &Config{
		ConfigManager: "127.0.0.1:2191",
		ProbeManager:  "127.0.0.1:2192",
		UserName:      "username",
		ApiKey:        "1234-1234-1234-1234",
		UpdateUrl:     "http://www.mytankstats.com/api.php?data=",
		UpdateMethod:  "Get",
		Probes: []probe.Probe{
			probe.Probe{"tank1", "i2c", "0x48", "ambient temp", "temp", 0},
			probe.Probe{"tank1", "1wire", "28-000005b3c94a", "aquarium temp", "temp", 0},
		},
	}
}
