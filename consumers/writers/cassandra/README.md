# Cassandra writer

Cassandra writer provides message repository implementation for Cassandra.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                             | Description                                                             | Default                        |
| ------------------------------------ | ----------------------------------------------------------------------- | ------------------------------ |
| SMQ_CASSANDRA_WRITER_LOG_LEVEL        | Log level for Cassandra writer (debug, info, warn, error)               | info                           |
| SMQ_CASSANDRA_WRITER_CONFIG_PATH      | Config file path with NATS subjects list, payload type and content-type | /config.toml                   |
| SMQ_CASSANDRA_WRITER_HTTP_HOST        | Cassandra service HTTP host                                             |                                |
| SMQ_CASSANDRA_WRITER_HTTP_PORT        | Cassandra service HTTP port                                             | 9004                           |
| SMQ_CASSANDRA_WRITER_HTTP_SERVER_CERT | Cassandra service HTTP server certificate path                          |                                |
| SMQ_CASSANDRA_WRITER_HTTP_SERVER_KEY  | Cassandra service HTTP server key path                                  |                                |
| SMQ_CASSANDRA_CLUSTER                 | Cassandra cluster comma separated addresses                             | 127.0.0.1                      |
| SMQ_CASSANDRA_KEYSPACE                | Cassandra keyspace name                                                 | supermq                     |
| SMQ_CASSANDRA_USER                    | Cassandra DB username                                                   | supermq                     |
| SMQ_CASSANDRA_PASS                    | Cassandra DB password                                                   | supermq                     |
| SMQ_CASSANDRA_PORT                    | Cassandra DB port                                                       | 9042                           |
| SMQ_MESSAGE_BROKER_URL                | Message broker instance URL                                             | nats://localhost:4222          |
| SMQ_JAEGER_URL                        | Jaeger server URL                                                       | http://jaeger:14268/api/traces |
| SMQ_SEND_TELEMETRY                    | Send telemetry to superMQ call home server                           | true                           |
| SMQ_CASSANDRA_WRITER_INSTANCE_ID      | Cassandra writer instance ID                                            |                                |

## Deployment

The service itself is distributed as Docker container. Check the [`cassandra-writer`](https://github.com/absmach/supermq-contrib/blob/main/docker/addons/cassandra-writer/docker-compose.yml#L30-L49) service section in docker-compose file to see how service is deployed.

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/absmach/supermq-contrib.git

cd supermq-contrib

# compile the cassandra writer
make cassandra-writer

# copy binary to bin
make install

# Set the environment variables and run the service
SMQ_CASSANDRA_WRITER_LOG_LEVEL=[Cassandra writer log level] \
SMQ_CASSANDRA_WRITER_CONFIG_PATH=[Config file path with NATS subjects list, payload type and content-type] \
SMQ_CASSANDRA_WRITER_HTTP_HOST=[Cassandra service HTTP host] \
SMQ_CASSANDRA_WRITER_HTTP_PORT=[Cassandra service HTTP port] \
SMQ_CASSANDRA_WRITER_HTTP_SERVER_CERT=[Cassandra service HTTP server cert] \
SMQ_CASSANDRA_WRITER_HTTP_SERVER_KEY=[Cassandra service HTTP server key] \
SMQ_CASSANDRA_CLUSTER=[Cassandra cluster comma separated addresses] \
SMQ_CASSANDRA_KEYSPACE=[Cassandra keyspace name] \
SMQ_CASSANDRA_USER=[Cassandra DB username] \
SMQ_CASSANDRA_PASS=[Cassandra DB password] \
SMQ_CASSANDRA_PORT=[Cassandra DB port] \
SMQ_MESSAGE_BROKER_URL=[Message Broker instance URL] \
SMQ_JAEGER_URL=[Jaeger server URL] \
SMQ_SEND_TELEMETRY=[Send telemetry to SuperMQ call home server] \
SMQ_CASSANDRA_WRITER_INSTANCE_ID=[Cassandra writer instance ID] \
$GOBIN/supermq-contrib-cassandra-writer
```

### Using docker-compose

This service can be deployed using docker containers. Docker compose file is
available in `<project_root>/docker/addons/cassandra-writer/docker-compose.yml`.
In order to run all SuperMQ core services, as well as mentioned optional ones,
execute following command:

```bash
./docker/addons/cassandra-writer/init.sh
```

## Usage

Starting service will start consuming normalized messages in SenML format.

[doc]: https://docs.supermq.abstractmachines.fr
