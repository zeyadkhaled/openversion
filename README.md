# Open Telemetry collector library

## Introduction

- This library offers a way to receive, process, and export all your data
(Metrics, Traces) in different formats and using mutliple services.

- There are two ways to run this library either as 

1)  An Agent that itself
    communicates with different receivers and then through a single channel
    exports to a collector and this collector then sends to multiple exporters.
    https://user-images.githubusercontent.com/10536136/48792454-2a69b900-eca9-11e8-96eb-c65b2b1e4e83.png
  
2)  The other way is by deploying a standalone collector service that is exposed
    though an OpenCensus reciever agent. This agent communicates over the
    network directly to multiple backend exporters.
    https://user-images.githubusercontent.com/10536136/46637070-65f05f80-cb0f-11e8-96e6-bc56468486b3.png

## Pipelines

- For either option above. The agent and the collector are structured in a YAML
  file that declares a service with pipelines.

The structure this file follows in this format:

```
receivers:
  opencensus:
    endpoint: "0.0.0.0:55678"

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
    namespace: promexample
logging:

processors:
  batch:
  queued_retry:

extensions:
  health_check:

service:
  extensions: [health_check]
  pipelines:
    traces:
      receivers: [opencensus]
      exporters: [logging]
      processors: [batch, queued_retry]
    metrics:
      receivers: [opencensus]
      exporters: [logging,prometheus]
```

The idea is to first declare your metrics or traces recievers.( In most cases it
is an opencensus receiver exposed on the endpoint of the collector or the
agent.)

Afterwards declare your Exporters then the processors and extensions which are
both optional.

At the end you declare a service with multiple pipelines. In the stucture above.
Give a name to the pipelines i.e traces, metrics, traces/2 and define their used
receivers and exporters at the least and then additionally processors.

https://github.com/open-telemetry/opentelemetry-collector/blob/master/docs/design.md


## Exporters

- There are multiple supported exporters for different purposes.
```
Supported trace exporters (sorted alphabetically):
    Jaeger
    OpenCensus
    OTLP
    Zipkin

Supported metric exporters (sorted alphabetically):
    OpenCensus
    Prometheus

Supported local exporters (sorted alphabetically):
    File
    Logging
```

Each exporter have different configuration but follows the same declaration language.
https://github.com/open-telemetry/opentelemetry-collector/tree/master/exporter

## Building the collector/agent

- For this step a docker-compose file is created.

- In this file declare all your used exporters and services. 

- Declare a collector and/or agent service which will expose a port and this
  will be used as an opencensus reciever.

- Include your application as service and make it depend on the collector or
  agent.

https://github.com/open-telemetry/opentelemetry-collector/blob/master/examples/demo/docker-compose.yaml


## Sending to openCensus agent all traces and metrics

- In your main.go of your application declare your OpenCensusAgent which is your
  agent or collectors recieving endpoint that is declared in docker-compose file.

```
ocAgentAddr, ok := os.LookupEnv("OTEL_AGENT_ENDPOINT")
	if !ok {
		ocAgentAddr = ocagent.DefaultAgentHost + ":" + string(ocagent.DefaultAgentPort)
	}
	oce, err := ocagent.NewExporter(
		ocagent.WithAddress(ocAgentAddr),
		ocagent.WithInsecure(),
		ocagent.WithServiceName(fmt.Sprintf("example-go-%d", os.Getpid())),
	)
```

- Register it for Metrics and Traces usage

```
trace.RegisterExporter(oce)
view.RegisterExporter(oce)
 ```

 - Check documentation of how to use this client library to send vendor
   agonistic metrics and traces that will be exported to all your registered
   exporters.
   go.opencensus.io/


## Running this 

- In the main directory where your docker-compose and otel-collector and otel-agent file are
  present run ``docker-compose up`` command.

- If every thing is successful you will start seeing your stats showing in your
  deployed exporters.

- If logging exporters with debug mode is added to your collector you will find
  output of all your collected metrics and traces in your docker output.