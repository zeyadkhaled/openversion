# Open Telemetry Collector with OTLP Exporter

## Introduction

The repository explains what is OpenTelemetry, OpenTelemetry Collector, OTLP
(Open Telemetry Protocol), and will proceed with a get started guideline to
integrate opentelemetry in your Go project and export your telemetry data to
different backends that are supported by the core Collector project and
additional backends coming from the community contributions.


### What is OpenTelemetry

OpenTelemetry provides a single set of APIs, libraries, agents, and collector
services to capture distributed traces and metrics from your application. 
[OpenTelemetry Website](https://opentelemetry.io/]

#### What are traces, metrics, and logs?

- Metrics: 
    - Everything from operating systems to applications generate metrics which,
      at the least, are going to include a name, a time stamp and a field to
      represent some value.

    - Most all metrics will enable you to tell if a resource is alive or dead,
      but if the target is valuable enough you’ll want to be able to ascertain
      what is actually wrong with the system or going wrong. 

- Traces: 
    - With so many application interdependencies these days, these operations
      will typically involve hops through multiple services (so called spans
      Traces, then, add critical visibility into the health of an application
      end-to-end. 

- Logs:
    - What’s more, logs tend to give more in-depth information about resources
      than metrics.  So, if metrics showed the resource is dead, logs will help
      tell you why it died.
      
[Explanation of Traces, Metrics, Logs](https://devops.com/metrics-logs-and-traces-the-golden-triangle-of-observability-in-monitoring/)
    

### What is OpenTelemetry Collector

- The OpenTelemetry Collector offers a vendor-agnostic implementation on how to receive, process, and export telemetry data. In addition, it removes the need to run, operate, and maintain multiple agents/collectors in order to support open-source telemetry data formats (e.g. Jaeger, Prometheus, etc.) sending to multiple open-source or commercial back-ends.
[https://github.com/open-telemetry/opentelemetry-collector](OpenTelemetry
Collector Repo)

- The name might be misleading that this collector offer directly exporting data
  in OpenTelemetry format; however, this is not the case but it is possible
  through OTLP exporters and receievers. 

### What is OTLP

- This is the protocol format that telemetry data could be collected in and then
  exported to a collector that has a receiver that understands this protocol
  format and could translate it to other backends like Prometheus, Jaeger, GCP.

### Why OpenTelemetry

- To date, two open-source projects have dominated the cloud-native telemetry landscape: OpenTracing and OpenCensus—each with its own telemetric standard and its own substantial community. In a commendable spirit of cooperation, however, the two projects have decided to converge in order to achieve an important mission: a single standard for built-in, high-quality cloud-native telemetry.

- Essentially OpenTelemetry converges the best of both projects into a single
  standard around which the entire cloud-native development community can rally.

[https://blog.thundra.io/what-is-opentelemetry-and-what-does-it-bring](Blog
Explaning What Exactly is OpenTelemetry)

### Why OpenTelemetry Collector

- Using just vanilla OpenTelemetry to export to multiple backends would require
  the usage of various exporters that are vendor specific which means more code,
  more dependencies, and high risk of breaking your project when OpenTelemetry
  updates their API but the vendor specific exporter is unchanged.
  (The reason this repo came to existance was a trial to integrate exporting
  OpenTelemetry data to Google Cloud Tracing and this ended in failure after the
  poorly maintained GCP OpenTelemetry exporter was dependent on old version of
  OpenTelemetry SDK)


## How to Get Started

### Understanding Open Telemetry Collector Architecture




