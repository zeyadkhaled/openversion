module opentelemetry.version.service

go 1.14

require (
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/go-redis/redis/v7 v7.3.0
	github.com/jackc/pgconn v1.5.0
	github.com/jackc/pgx/v4 v4.6.0
	github.com/rs/cors v1.7.0
	github.com/rs/zerolog v1.18.0
	go.opentelemetry.io/otel v0.8.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.8.0
	go.opentelemetry.io/otel/exporters/otlp v0.8.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/net v0.0.0-20200501053045-e0ff5e5a1de5 // indirect
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)
