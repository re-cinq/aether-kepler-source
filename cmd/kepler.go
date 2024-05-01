package main

import (
	"fmt"
	"log/slog"
	"os"
	"v1/pkg/config"
	"v1/pkg/kepler"

	"github.com/hashicorp/go-plugin"
	"github.com/prometheus/client_golang/api"
	apiconfig "github.com/prometheus/common/config"

	aetherplugin "github.com/re-cinq/aether/pkg/plugin"
)

func main() {

	// in order to log from plugin you have to output to Stderr
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg := &config.Config{}

	err := cfg.Load()
	if err != nil {
		logger.Error("failed loading config", "error", err)
		os.Exit(1)
	}

	client, err := api.NewClient(api.Config{
		Address: fmt.Sprintf("%v:%v", cfg.PrometheusURL, cfg.PrometheusPort),
		// TODO: setup authentication
		RoundTripper: apiconfig.NewBasicAuthRoundTripper("admin", "", "", "", api.DefaultRoundTripper),
	})
	if err != nil {
		logger.Error("unable to setup prometheus client", "error", err)
		os.Exit(-1)
	}

	// initialize the struct
	k := kepler.New(
		kepler.WithLogger(logger),
		kepler.WithConfig(cfg),
		kepler.WithPrometheusClient(&client),
	)

	// The name of the plugin needs to match the name of the binary
	// created by the build process.
	pluginMap := map[string]plugin.Plugin{
		"kepler": &aetherplugin.SourcePlugin{Impl: k},
	}

	// This is a blocking call that will start the plugin
	plugin.Serve(&plugin.ServeConfig{
		// This should corrspond to the handshake that the version of
		// aether you are using is using
		HandshakeConfig: aetherplugin.SourceHandshake,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
