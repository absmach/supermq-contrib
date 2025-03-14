# MongoDB reader

MongoDB reader provides message repository implementation for MongoDB.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                         | Description                                         | Default                        |
| -------------------------------- | --------------------------------------------------- | ------------------------------ |
| SMQ_MONGO_READER_LOG_LEVEL        | Service log level                                   | info                           |
| SMQ_MONGO_READER_HTTP_HOST        | Service HTTP host                                   | localhost                      |
| SMQ_MONGO_READER_HTTP_PORT        | Service HTTP port                                   | 9007                           |
| SMQ_MONGO_READER_HTTP_SERVER_CERT | Service HTTP server cert                            | ""                             |
| SMQ_MONGO_READER_HTTP_SERVER_KEY  | Service HTTP server key                             | ""                             |
| SMQ_MONGO_NAME                    | MongoDB database name                               | messages                       |
| SMQ_MONGO_HOST                    | MongoDB database host                               | localhost                      |
| SMQ_MONGO_PORT                    | MongoDB database port                               | 27017                          |
| SMQ_CLIENTS_AUTH_GRPC_URL          | Clients service Auth gRPC URL                        | localhost:7000                 |
| SMQ_CLIENTS_AUTH_GRPC_TIMEOUT      | Clients service Auth gRPC request timeout in seconds | 1s                             |
| SMQ_CLIENTS_AUTH_GRPC_CLIENT_TLS   | Flag that indicates if TLS should be turned on      | false                          |
| SMQ_CLIENTS_AUTH_GRPC_CA_CERTS     | Path to trusted CAs in PEM format                   | ""                             |
| SMQ_AUTH_GRPC_URL                 | Auth service gRPC URL                               | localhost:7001                 |
| SMQ_AUTH_GRPC_TIMEOUT             | Auth service gRPC request timeout in seconds        | 1s                             |
| SMQ_AUTH_GRPC_CLIENT_TLS          | Flag that indicates if TLS should be turned on      | false                          |
| SMQ_AUTH_GRPC_CA_CERT             | Path to trusted CAs in PEM format                   | ""                             |
| SMQ_JAEGER_URL                    | Jaeger server URL                                   | http://jaeger:14268/api/traces |
| SMQ_SEND_TELEMETRY                | Send telemetry to supermq call home server       | true                           |
| SMQ_MONGO_READER_INSTANCE_ID      | Service instance ID                                 | ""                             |

## Deployment

The service itself is distributed as Docker container. Check the [`mongodb-reader`](https://github.com/absmach/supermq-contrib/blob/main/docker/addons/mongodb-reader/docker-compose.yml#L16-L37) service section in
docker-compose file to see how service is deployed.

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/absmach/supermq-contrib.git

cd supermq-contrib

# compile the mongodb reader
make mongodb-reader

# copy binary to bin
make install

# Set the environment variables and run the service
SMQ_MONGO_READER_LOG_LEVEL=[Service log level] \
SMQ_MONGO_READER_HTTP_HOST=[Service HTTP host] \
SMQ_MONGO_READER_HTTP_PORT=[Service HTTP port] \
SMQ_MONGO_READER_HTTP_SERVER_CERT=[Path to server pem certificate file] \
SMQ_MONGO_READER_HTTP_SERVER_KEY=[Path to server pem key file] \
SMQ_MONGO_NAME=[MongoDB database name] \
SMQ_MONGO_HOST=[MongoDB database host] \
SMQ_MONGO_PORT=[MongoDB database port] \
SMQ_CLIENTS_AUTH_GRPC_URL=[Clients service Auth gRPC URL] \
SMQ_CLIENTS_AUTH_GRPC_TIMEOUT=[Clients service Auth gRPC request timeout in seconds] \
SMQ_CLIENTS_AUTH_GRPC_CLIENT_TLS=[Flag that indicates if TLS should be turned on] \
SMQ_CLIENTS_AUTH_GRPC_CA_CERTS=[Path to trusted CAs in PEM format] \
SMQ_AUTH_GRPC_URL=[Auth service gRPC URL] \
SMQ_AUTH_GRPC_TIMEOUT=[Auth service gRPC request timeout in seconds] \
SMQ_AUTH_GRPC_CLIENT_TLS=[Flag that indicates if TLS should be turned on] \
SMQ_AUTH_GRPC_CA_CERT=[Path to trusted CAs in PEM format] \
SMQ_JAEGER_URL=[Jaeger server URL] \
SMQ_SEND_TELEMETRY=[Send telemetry to supermq call home server] \
SMQ_MONGO_READER_INSTANCE_ID=[Service instance ID] \
$GOBIN/supermq-contrib-mongodb-reader

```

### Using docker-compose

This service can be deployed using docker containers. Docker compose file is
available in `<project_root>/docker/addons/mongodb-reader/docker-compose.yml`.
In order to run all SupeMQ core services, as well as mentioned optional ones,
execute following command:

```bash
docker compose -f docker/docker-compose.yml up -d
docker compose -f docker/addons/mongodb-reader/docker-compose.yml up -d
```

## Usage

Service exposes [HTTP API](https://docs.api.supermq.abstractmachines.fr/?urls.primaryName=readers-openapi.yml) for fetching messages.

```
Note: MongoDB Reader doesn't support searching substrings from string_value, due to inefficient searching as the current data model is not suitable for this type of queries.
```

[doc]: https://docs.supermq.abstractmachines.fr
