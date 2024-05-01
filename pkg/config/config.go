package config

import (
	"time"

	"github.com/caarlos0/env/v10"
)

// Config is used to fetch configuration from the environment
// due to how the plugin is run (alongside Aether) we dont really have access
// to Aethers configuration so the easiest way to do this is to just use
// environments for now. However this might change in the future
type Config struct {
	// Interval is the window of aggregated data scraped from the source
	// This value **NEEDS** to match the interval configuration option
	// passed into the aether local.config
	Interval time.Duration `env:"INTERVAL,required"`

	// Provider is the cloud provider instance data is being scraped from
	// Currently this value can only be "aws" or "gcp"
	Provider string `env:"PROVIDER,required"`

	// PrometheusURL is the URL to the Prometheus server
	// where kepler metrics will be collected from
	PrometheusURL string `env:"PROMETHEUS_URL,required"`

	// PrometheusPort is the port Prometheus is running on
	// By default this is set to 9090
	PrometheusPort string `env:"PROMETHEUS_PORT" envDefault:"9090"`

	// TODO: Set up authentication with prometheus
	//AuthToken string        `env:"AUTH_TOKEN"`
}

// Load parses the environment into our Config Struct
func (c *Config) Load() error {
	return env.Parse(c)
}
