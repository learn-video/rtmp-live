package discovery

import "github.com/caarlos0/env/v9"

type Config struct {
	HLSPath         string `env:"HLS_PATH"`
	IP              string `env:"IP"`
	DiscoveryAPIURL string `env:"DISCOVERY_API_URL"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
