# InfluxDB reader

InfluxDB reader provides message repository implementation for InfluxDB.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                         | Description                                         | Default                        |
| -------------------------------- | --------------------------------------------------- | ------------------------------ |
| SMQ_INFLUX_READER_LOG_LEVEL       | Service log level                                   | info                           |
| SMQ_INFLUX_READER_HTTP_HOST       | Service HTTP host                                   | localhost                      |
| SMQ_INFLUX_READER_HTTP_PORT       | Service HTTP port                                   | 9005                           |
| SMQ_INFLUX_READER_SERVER_CERT     | Service HTTP server cert                            | ""                             |
| SMQ_INFLUX_READER_SERVER_KEY      | Service HTTP server key                             | ""                             |
| SMQ_INFLUXDB_PROTOCOL             | InfluxDB protocol                                   | http                           |
| SMQ_INFLUXDB_HOST                 | InfluxDB host name                                  | localhost                      |
| SMQ_INFLUXDB_PORT                 | Default port of InfluxDB database                   | 8086                           |
| SMQ_INFLUXDB_ADMIN_USER           | Default user of InfluxDB database                   | supermq                     |
| SMQ_INFLUXDB_ADMIN_PASSWORD       | Default password of InfluxDB user                   | supermq                     |
| SMQ_INFLUXDB_NAME                 | InfluxDB database name                              | supermq                     |
| SMQ_INFLUXDB_BUCKET               | InfluxDB bucket name                                | supermq-bucket              |
| SMQ_INFLUXDB_ORG                  | InfluxDB organization name                          | supermq                     |
| SMQ_INFLUXDB_TOKEN                | InfluxDB API token                                  | supermq-token               |
| SMQ_INFLUXDB_DBURL                | InfluxDB database URL                               | ""                             |
| SMQ_INFLUXDB_USER_AGENT           | InfluxDB user agent                                 | ""                             |
| SMQ_INFLUXDB_TIMEOUT              | InfluxDB client connection readiness timeout        | 1s                             |
| SMQ_INFLUXDB_INSECURE_SKIP_VERIFY | InfluxDB insecure skip verify                       | false                          |
| SMQ_CLIENTS_AUTH_GRPC_URL          | Clients service Auth gRPC URL                        | localhost:7000                 |
| SMQ_CLIENTS_AUTH_GRPC_TIMEOUT      | Clients service Auth gRPC request timeout in seconds | 1s                             |
| SMQ_CLIENTS_AUTH_GRPC_CLIENT_TLS   | Flag that indicates if TLS should be turned on      | false                          |
| SMQ_CLIENTS_AUTH_GRPC_CA_CERTS     | Path to trusted CAs in PEM format                   | ""                             |
| SMQ_AUTH_GRPC_URL                 | Auth service gRPC URL                               | localhost:7001                 |
| SMQ_AUTH_GRPC_TIMEOUT             | Auth service gRPC request timeout in seconds        | 1s                             |
| SMQ_AUTH_GRPC_CLIENT_TLS          | Flag that indicates if TLS should be turned on      | false                          |
| SMQ_AUTH_GRPC_CA_CERTS            | Path to trusted CAs in PEM format                   | ""                             |
| SMQ_JAEGER_URL                    | Jaeger server URL                                   | http://jaeger:14268/api/traces |
| SMQ_SEND_TELEMETRY                | Send telemetry to supermq call home server       | true                           |
| SMQ_INFLUX_READER_INSTANCE_ID     | InfluxDB reader instance ID                         |                                |

## Deployment

The service itself is distributed as Docker container. Check the [`influxdb-reader`](https://github.com/absmach/supermq-contrib/blob/main/docker/addons/influxdb-reader/docker-compose.yml#L17-L40) service section in docker-compose file to see how service is deployed.

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/absmach/supermq-contrib.git

cd supermq-contrib

# compile the influxdb-reader
make influxdb-reader

# copy binary to bin
make install

# Set the environment variables and run the service
SMQ_INFLUX_READER_LOG_LEVEL=[Service log level] \
SMQ_INFLUX_READER_HTTP_HOST=[Service HTTP host] \
SMQ_INFLUX_READER_HTTP_PORT=[Service HTTP port] \
SMQ_INFLUX_READER_HTTP_SERVER_CERT=[Service HTTP server certificate] \
SMQ_INFLUX_READER_HTTP_SERVER_KEY=[Service HTTP server key] \
SMQ_INFLUXDB_PROTOCOL=[InfluxDB protocol] \
SMQ_INFLUXDB_HOST=[InfluxDB database host] \
SMQ_INFLUXDB_PORT=[InfluxDB database port] \
SMQ_INFLUXDB_ADMIN_USER=[InfluxDB admin user] \
SMQ_INFLUXDB_ADMIN_PASSWORD=[InfluxDB admin password] \
SMQ_INFLUXDB_NAME=[InfluxDB database name] \
SMQ_INFLUXDB_BUCKET=[InfluxDB bucket] \
SMQ_INFLUXDB_ORG=[InfluxDB org] \
SMQ_INFLUXDB_TOKEN=[InfluxDB token] \
SMQ_INFLUXDB_DBURL=[InfluxDB database URL] \
SMQ_INFLUXDB_USER_AGENT=[InfluxDB user agent] \
SMQ_INFLUXDB_TIMEOUT=[InfluxDB timeout] \
SMQ_INFLUXDB_INSECURE_SKIP_VERIFY=[InfluxDB insecure skip verify] \
SMQ_CLIENTS_AUTH_GRPC_URL=[Clients service Auth gRPC URL] \
SMQ_CLIENTS_AURH_GRPC_TIMEOUT=[Clients service Auth gRPC request timeout in seconds] \
SMQ_CLIENTS_AUTH_GRPC_CLIENT_TLS=[Flag that indicates if TLS should be turned on] \
SMQ_CLIENTS_AUTH_GRPC_CA_CERTS=[Path to trusted CAs in PEM format] \
SMQ_AUTH_GRPC_URL=[Auth service gRPC URL] \
SMQ_AUTH_GRPC_TIMEOUT=[Auth service gRPC request timeout in seconds] \
SMQ_AUTH_GRPC_CLIENT_TLS=[Flag that indicates if TLS should be turned on] \
SMQ_AUTH_GRPC_CA_CERTS=[Path to trusted CAs in PEM format] \
SMQ_JAEGER_URL=[Jaeger server URL] \
SMQ_SEND_TELEMETRY=[Send telemetry to supermq call home server] \
SMQ_INFLUX_READER_INSTANCE_ID=[InfluxDB reader instance ID] \
$GOBIN/supermq-contrib-influxdb

```

### Using docker-compose

This service can be deployed using docker containers. Docker compose file is
available in `<project_root>/docker/addons/influxdb-reader/docker-compose.yml`.
In order to run all SupeMQ core services, as well as mentioned optional ones,
execute following command:

```bash
docker compose -f docker/docker-compose.yml up -d
docker compose -f docker/addons/influxdb-reader/docker-compose.yml up -d
```

And, to use the default .env file, execute the following command:

```bash
docker compose -f docker/addons/influxdb-reader/docker-compose.yml up --env-file docker/.env -d
```

## Usage

Service exposes [HTTP API](https://docs.api.supermq.abstractmachines.fr/?urls.primaryName=readers-openapi.yml) for fetching messages.

Comparator Usage Guide:
| Comparator | Usage | Example |  
|----------------------|-----------------------------------------------------------------------------|------------------------------------|
| eq | Return values that are equal to the query | eq["active"] -> "active" |  
| ge | Return values that are substrings of the query | ge["tiv"] -> "active" and "tiv" |  
| gt | Return values that are substrings of the query and not equal to the query | gt["tiv"] -> "active" |  
| le | Return values that are superstrings of the query | le["active"] -> "tiv" |  
| lt | Return values that are superstrings of the query and not equal to the query | lt["active"] -> "active" and "tiv" |

Official docs can be found [here](https://docs.supermq.abstractmachines.fr).
