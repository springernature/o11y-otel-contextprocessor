# Basic OpenTelemetry "hello-world" processor

The purpose of this code is to show how easy is to integrate a new component in OpenTelemetry collector.
This is a *hello-world processor*, which prints all resource attributes associated with metrics and if there a match with the configuration parameter name, prints its value.

# About

This is a quick hands-on guide, describing a basic and minimum setup of a custom metrics processor. There are no tests, information about styles, best practices ... You can see this repo as a short version of https://opentelemetry.io/docs/collector/custom-collector/

Also, to get starting contributing in OpenTelemetry, please have a look at these resources:

* OpenTelemetry Collector: https://github.com/open-telemetry/opentelemetry-collector/blob/main/CONTRIBUTING.md
* OpenTelemetry Collector-Contrib: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/CONTRIBUTING.md

# Starting ...

1. Get golang installed: https://go.dev/doc/install
2. Install the [ocb](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder) tool: https://github.com/open-telemetry/opentelemetry-collector/releases

## ... from zero

If you want to get to your own custom

1. Generate a `builder-config.yaml`. Adjust `name`, `output_path`, `otelcol_version` and the gomod versions, if needed.
    ```
    cat <<EOF > builder-config.yaml
    dist:
        name: otelcol-de
        description: Basic OTel Collector distribution for Developers
        output_path: ./otelcol-dev
        otelcol_version: 0.88.0

    extensions:
    - gomod: go.opentelemetry.io/collector/extension/zpagesextension v0.88.0
    - gomod: go.opentelemetry.io/collector/extension/ballastextension v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.88.0

    exporters:
    - gomod: go.opentelemetry.io/collector/exporter/debugexporter v0.88.0
    - gomod: go.opentelemetry.io/collector/exporter/otlpexporter v0.88.0
    - gomod: go.opentelemetry.io/collector/exporter/otlphttpexporter v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.88.0

    processors:
    - gomod: go.opentelemetry.io/collector/processor/batchprocessor v0.88.0
    - gomod: go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.88.0

    receivers:
    - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.88.0
    - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.88.0

    connectors:
    - gomod: go.opentelemetry.io/collector/connector/forwardconnector v0.88.0
    EOF
    ```
2. Run `ocb --config builder-config.yaml`. A new directory `otelcol-dev` (or what you have defined in `output_path`) is created. The code generated should not be touched, all changes should be done via the `builder-config.yaml` file. You can see the builder configuration of the OpenTelemetry contrib here: https://github.com/open-telemetry/opentelemetry-collector/blob/main/cmd/otelcorecol/builder-config.yaml
3. Go to new `otelcold-dev` and you will see a new executable `otelcol-dev` (or the value defined in `name`). You can execute it by providing a basic configuration file `otelcold-dev --config otelcol.yaml`:
    ```
    cat <<EOF > otelcol.yaml
    extensions:
    health_check:
    pprof:
        endpoint: 0.0.0.0:1777
    zpages:
        endpoint: 0.0.0.0:55679

    receivers:
    prometheus:
        config:
        scrape_configs:
            - job_name: 'otel-collector'
            scrape_interval: 30s
            static_configs:
                - targets: ['0.0.0.0:8888']

    exporters:
    debug:
        verbosity: detailed

    service:
    telemetry:
        logs:
        level: "debug"
    pipelines:
        metrics:
        receivers: [prometheus]
        processors: []
        exporters: [debug]
    extensions: 
    - health_check
    - pprof
    - zpages
    ```
    With this configuration the collector will scrape its own metrics.

## ... from the current repo

You can have a look to the folder `otelcol-dev/helloworldmetricsprocessor` and start modifying the code. Or

1. Create a new directory like `helloworldmetricsprocessor`.
2. Create a new golang module `go mod init github.com/jriguera/otel-helloworldprocessor/otelcol-dev/helloworldmetricsprocessor` (change according to your choice)
3. Create new golang files. Use same package with same name as the directory `helloworldmetricsprocessor`:
   1. `config.go` to define the component parameters and the function to validate them.
   2. `factory.go` to define the facture which create instances of the new component. Use the name `NewFactory` and the factory provided by the component package (eg `processor.NewFactory`)
   3. Other files to run the task (eg `processor.go`).
4. Add the component to the section of `builder-config.yaml` with the same reference you used in the `go mod init ...`. Create another section to replace this module name with the local directory created in point 1.
   ```
   processors:
   - gomod: github.com/jriguera/otel-helloworldprocessor/otelcol-dev/helloworldmetricsprocessor v0.88.0

   replaces:
   - github.com/jriguera/otel-helloworldprocessor/otelcol-dev/helloworldmetricsprocessor => ./helloworldmetricsprocessor
   ```
5. Finally run `ocb --config builder-config.yaml` to generate the binary with the new component

