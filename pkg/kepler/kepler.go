package kepler

import (
	"log/slog"
	"v1/pkg/config"

	"github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"

	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

type KeplerSource struct {
	logger *slog.Logger

	cfg    *config.Config
	client *api.Client
	v1API  prometheus.API

	instancesMap map[string]*v1.Instance
}

type option func(*KeplerSource)

func WithLogger(l *slog.Logger) option {
	return func(k *KeplerSource) {
		k.logger = l
	}
}

func WithConfig(c *config.Config) option {
	return func(k *KeplerSource) {
		k.cfg = c
	}
}

func WithPrometheusClient(c *api.Client) option {
	return func(k *KeplerSource) {
		k.client = c
	}
}

func New(opts ...option) *KeplerSource {
	// set defaults
	r := &KeplerSource{
		instancesMap: make(map[string]*v1.Instance),
	}

	// overwrite options
	for _, o := range opts {
		o(r)
	}

	return r
}
