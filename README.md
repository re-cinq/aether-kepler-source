# aether-kelper-source
The plugin is designed as a source plugin for the aether project. It fetches data from the kepler API and sends it to the aether project via gRPC calls and the go-plugin framework.

### Installation
To install the plugin with aether in cluster environments follow the [installation guide][1].

### Configuration
The following environment variables are required to to run the plugin:

```
INTERVAL=<interval>
PROVIDER=<provider>
PROMETHEUS_URL=<prometheus_url>

# optional
PROMETHEUS_PORT=<prometheus_port>
```

`<interval>` is the interval (5m/30m) at which the plugin should fetch data from the kepler API. This value **NEEDS** to match that of the `scrapingInterval` value specified in the `local.yaml` file in the `aether` project.

`<provider>` is the Cloud Provider of what kepler is collecting metrics on. At the moment the only supported values are `aws` and `gcp`. With future support for `azure` and `on-prem`.

`<prometheus_url>` is the URL of the prometheus server that the plugin scrapes kepler metrics from.


`<prometheus_port>` is the port of the prometheus server that the plugin scrapes kepler metrics from. This is an optional field and defaults to `9090`.


### Running locally with aether

Aether runs and is managed with the `Dockerfile.dev`, `docker-compose.yml`, and `local.yaml` files.

To get the plugin running within aether:
1. Build the plugin and store it in the `.plugin` directory

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o .plugin/kepler cmd/kepler.go
```
__NOTE__ If you want to run the plugin stand-alone, although not it's design, you can do that with:
```bash
./.plugin/kepler
```
2. Specify the "source" directory that aether will read the plugin binary from. This is a configuration in the `local.yaml` file.

```yaml
plugins:
  SourceDir: .plugins/
```
3. An `.env` file is required to be present in the root directory of the project.
This file should contain the environment variables specified in the configuration section above.

```
4. Run aether using docker-compose:
```
docker compose up
```

[1]: https://aether.green/docs/tutorials/kepler/
