# Cassandra reader

Cassandra reader provides message repository implementation for Cassandra.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                             | Description                                         | Default                        |
| ------------------------------------ | --------------------------------------------------- | ------------------------------ |
| SMQ_CASSANDRA_READER_LOG_LEVEL        | Cassandra service log level                         | debug                          |
| SMQ_CASSANDRA_READER_HTTP_HOST        | Cassandra service HTTP host                         | localhost                      |
| SMQ_CASSANDRA_READER_HTTP_PORT        | Cassandra service HTTP port                         | 9003                           |
| SMQ_CASSANDRA_READER_HTTP_SERVER_CERT | Cassandra service HTTP server cert                  | ""                             |
| SMQ_CASSANDRA_READER_HTTP_SERVER_KEY  | Cassandra service HTTP server key                   | ""                             |
| SMQ_CASSANDRA_CLUSTER                 | Cassandra cluster comma separated addresses         | localhost                      |
| SMQ_CASSANDRA_USER                    | Cassandra DB username                               | supermq                     |
| SMQ_CASSANDRA_PASS                    | Cassandra DB password                               | supermq                     |
| SMQ_CASSANDRA_KEYSPACE                | Cassandra keyspace name                             | messages                       |
| SMQ_CASSANDRA_PORT                    | Cassandra DB port                                   | 9042                           |
| SMQ_CLIENTS_AUTH_GRPC_URL              | Clients service Auth gRPC URL                        | localhost:7000                 |
| SMQ_CLIENTS_AUTH_GRPC_TIMEOUT          | Clients service Auth gRPC request timeout in seconds | 1                              |
| SMQ_CLIENTS_AUTH_GRPC_CLIENT_TLS       | Clients service Auth gRPC TLS enabled                | false                          |
| SMQ_CLIENTS_AUTH_GRPC_CA_CERTS         | Clients service Auth gRPC CA certificates            | ""                             |
| SMQ_AUTH_GRPC_URL                     | Auth service gRPC URL                               | localhost:7001                 |
| SMQ_AUTH_GRPC_TIMEOUT                 | Auth service gRPC request timeout in seconds        | 1s                             |
| SMQ_AUTH_GRPC_CLIENT_TLS              | Auth service gRPC TLS enabled                       | false                          |
| SMQ_AUTH_GRPC_CA_CERT                 | Auth service gRPC CA certificates                   | ""                             |
| SMQ_JAEGER_URL                        | Jaeger server URL                                   | http://jaeger:14268/api/traces |
| SMQ_SEND_TELEMETRY                    | Send telemetry to supermq call home server       | true                           |
| SMQ_CASSANDRA_READER_INSTANCE_ID      | Cassandra Reader instance ID                        | ""                             |

## Deployment

The service itself is distributed as Docker container. Check the [`cassandra-reader`](https://github.com/absmach/supermq-contrib/blob/main/docker/addons/cassandra-reader/docker-compose.yml#L15-L35) service section in
docker-compose file to see how service is deployed.

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/absmach/supermq-contrib.git

cd supermq-contrib

# compile the cassandra
make cassandra-reader

# copy binary to bin
make install

# Set the environment variables and run the service
SMQ_CASSANDRA_READER_LOG_LEVEL=[Cassandra Service log level] \
SMQ_CASSANDRA_READER_HTTP_HOST=[Cassandra Service HTTP host] \
SMQ_CASSANDRA_READER_HTTP_PORT=[Cassandra Service HTTP port] \
SMQ_CASSANDRA_READER_HTTP_SERVER_CERT=[Cassandra Service HTTP server cert] \
SMQ_CASSANDRA_READER_HTTP_SERVER_KEY=[Cassandra Service HTTP server key] \
SMQ_CASSANDRA_CLUSTER=[Cassandra cluster comma separated addresses] \
SMQ_CASSANDRA_KEYSPACE=[Cassandra keyspace name] \
SMQ_CASSANDRA_USER=[Cassandra DB username] \
SMQ_CASSANDRA_PASS=[Cassandra DB password] \
SMQ_CASSANDRA_PORT=[Cassandra DB port] \
SMQ_CLIENTS_AUTH_GRPC_URL=[Clients service Auth gRPC URL] \
SMQ_CLIENTS_AUTH_GRPC_TIMEOUT=[Clients service Auth gRPC request timeout in seconds] \
SMQ_CLIENTS_AUTH_GRPC_CLIENT_TLS=[Clients service Auth gRPC TLS enabled] \
SMQ_CLIENTS_AUTH_GRPC_CA_CERTS=[Clients service Auth gRPC CA certificates] \
SMQ_AUTH_GRPC_URL=[Auth service gRPC URL] \
SMQ_AUTH_GRPC_TIMEOUT=[Auth service gRPC request timeout in seconds] \
SMQ_AUTH_GRPC_CLIENT_TLS=[Auth service gRPC TLS enabled] \
SMQ_AUTH_GRPC_CA_CERT=[Auth service gRPC CA certificates] \
SMQ_JAEGER_URL=[Jaeger server URL] \
SMQ_SEND_TELEMETRY=[Send telemetry to supermq call home server] \
SMQ_CASSANDRA_READER_INSTANCE_ID=[Cassandra Reader instance ID] \
$GOBIN/supermq-contrib-cassandra-reader
```

### Using docker-compose

This service can be deployed using docker containers. Docker compose file is
available in `<project_root>/docker/addons/cassandra-reader/docker-compose.yml`.
In order to run all SupeMQ core services, as well as mentioned optional ones,
execute following command:

```bash
docker compose -f docker/docker-compose.yml up -d
./docker/addons/cassandra-writer/init.sh
docker compose -f docker/addons/casandra-reader/docker-compose.yml up -d
```

## Usage

Service exposes [HTTP API](https://docs.api.supermq.abstractmachines.fr/?urls.primaryName=readers-openapi.yml) for fetching messages.

```
Note: Cassandra Reader doesn't support searching substrings from string_value, due to inefficient searching as the current data model is not suitable for this type of queries.
```

[doc]: https://docs.supermq.abstractmachines.fr
