package config

import (
	"regexp"

	"github.com/BurntSushi/toml"
)

type RewriteRule struct {
	Match   string `toml:"match"`
	Replace string `toml:"replace"`
	Pattern *regexp.Regexp
}

type Route struct {
	Path        string        `toml:"path"`
	Backend     string        `toml:"backend"`
	StripPrefix bool          `toml:"strip_prefix"`
	Rewrites    []RewriteRule `toml:"rewrites"`
}

type GatewayConfig struct {
	Port   int     `toml:"port"`
	Routes []Route `toml:"routes"`
}

func LoadConfig(filePath string) (*GatewayConfig, error) {
	var cfg GatewayConfig
	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return nil, err
	}

	for i := range cfg.Routes {
		for j := range cfg.Routes[i].Rewrites {
			compiled, err := regexp.Compile(cfg.Routes[i].Rewrites[j].Match)
			if err != nil {
				return nil, err
			}
			cfg.Routes[i].Rewrites[j].Pattern = compiled
		}
	}

	return &cfg, nil
}
