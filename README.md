# aether-kelper-source
The plugin is designed as a source plugin for the aether project. It fetches data from the kepler API and sends it to the aether project via gRPC calls and the go-plugin framework.

## Query and calculations
Currently we query for the Kepler CPU and memory energy consumptions with the following queries:
```
sum without(mode) (rate(kepler_container_core_joules_total[5m])
sum without(mode) (rate(kepler_container_dram_joules_total[5m])
```

The `5m` is configurable based on the `INTERVAL` environment variable, however we strongly recommend that the `INTERVAL` value in the plugin matches the `scrapingInterval` value in the `local.yaml` file in the `aether` project.

The `sum without (mode)` portion of the query is to remove the `mode` label from the query by summing it's two parts (idle and dynamic) to get the absolute power value.
To be more specific, from a [Kepler Blog][2]: "The dynamic power is directly related to the resource utilization and the idle power is the constant power that does not vary regardless if the system is at rest or with load."

The result from the query is thus the absolute energy consumption in Joules summed over the `INTERVAL` period.
The aether project expects the result to be in kilowatt hours, so we convert the result to kilowatt hours by dividing the result by the interval time in seconds, then by `3600000` (to convert from Joules/sec).

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
[2]: https://www.cncf.io/blog/2023/10/11/exploring-keplers-potentials-unveiling-cloud-application-power-consumption/
