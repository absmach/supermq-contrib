# InfluxDB writer

InfluxDB writer provides message repository implementation for InfluxDB.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                          | Description                                                                       | Default                        |
| --------------------------------- | --------------------------------------------------------------------------------- | ------------------------------ |
| SMQ_INFLUX_WRITER_LOG_LEVEL        | Log level for InfluxDB writer (debug, info, warn, error)                          | info                           |
| SMQ_INFLUX_WRITER_CONFIG_PATH      | Config file path with message broker subjects list, payload type and content-type | /configs.toml                  |
| SMQ_INFLUX_WRITER_HTTP_HOST        | Service HTTP host                                                                 |                                |
| SMQ_INFLUX_WRITER_HTTP_PORT        | Service HTTP port                                                                 | 9006                           |
| SMQ_INFLUX_WRITER_HTTP_SERVER_CERT | Path to server certificate in pem format                                          |                                |
| SMQ_INFLUX_WRITER_HTTP_SERVER_KEY  | Path to server key in pem format                                                  |                                |
| SMQ_INFLUXDB_PROTOCOL              | InfluxDB protocol                                                                 | http                           |
| SMQ_INFLUXDB_HOST                  | InfluxDB host name                                                                | supermq-influxdb            |
| SMQ_INFLUXDB_PORT                  | Default port of InfluxDB database                                                 | 8086                           |
| SMQ_INFLUXDB_ADMIN_USER            | Default user of InfluxDB database                                                 | supermq                     |
| SMQ_INFLUXDB_ADMIN_PASSWORD        | Default password of InfluxDB user                                                 | supermq                     |
| SMQ_INFLUXDB_NAME                  | InfluxDB database name                                                            | supermq                     |
| SMQ_INFLUXDB_BUCKET                | InfluxDB bucket name                                                              | supermq-bucket              |
| SMQ_INFLUXDB_ORG                   | InfluxDB organization name                                                        | supermq                     |
| SMQ_INFLUXDB_TOKEN                 | InfluxDB API token                                                                | supermq-token               |
| SMQ_INFLUXDB_DBURL                 | InfluxDB database URL                                                             |                                |
| SMQ_INFLUXDB_USER_AGENT            | InfluxDB user agent                                                               |                                |
| SMQ_INFLUXDB_TIMEOUT               | InfluxDB client connection readiness timeout                                      | 1s                             |
| SMQ_INFLUXDB_INSECURE_SKIP_VERIFY  | InfluxDB client connection insecure skip verify                                   | false                          |
| SMQ_MESSAGE_BROKER_URL             | Message broker instance URL                                                       | nats://localhost:4222          |
| SMQ_JAEGER_URL                     | Jaeger server URL                                                                 | http://jaeger:14268/api/traces |
| SMQ_SEND_TELEMETRY                 | Send telemetry to supermq call home server                                     | true                           |
| SMQ_INFLUX_WRITER_INSTANCE_ID      | InfluxDB writer instance ID                                                       |                                |

## Deployment

The service itself is distributed as Docker container. Check the [`influxdb-writer`](https://github.com/absmach/supermq-contrib/blob/main/docker/addons/influxdb-writer/docker-compose.yml#L35-L58) service section in docker-compose file to see how service is deployed.

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/absmach/supermq-contrib.git

cd supermq-contrib

# compile the influxdb
make influxdb

# copy binary to bin
make install

# Set the environment variables and run the service
SMQ_INFLUX_WRITER_LOG_LEVEL=[Influx writer log level] \
SMQ_INFLUX_WRITER_CONFIG_PATH=[Config file path with Message broker subjects list, payload type and content-type] \
SMQ_INFLUX_WRITER_HTTP_HOST=[Service HTTP host] \
SMQ_INFLUX_WRITER_HTTP_PORT=[Service HTTP port] \
SMQ_INFLUX_WRITER_HTTP_SERVER_CERT=[Service HTTP server cert] \
SMQ_INFLUX_WRITER_HTTP_SERVER_KEY=[Service HTTP server key] \
SMQ_INFLUXDB_PROTOCOL=[InfluxDB protocol] \
SMQ_INFLUXDB_HOST=[InfluxDB database host] \
SMQ_INFLUXDB_PORT=[InfluxDB database port] \
SMQ_INFLUXDB_ADMIN_USER=[InfluxDB admin user] \
SMQ_INFLUXDB_ADMIN_PASSWORD=[InfluxDB admin password] \
SMQ_INFLUXDB_NAME=[InfluxDB database name] \
SMQ_INFLUXDB_BUCKET=[InfluxDB bucket] \
SMQ_INFLUXDB_ORG=[InfluxDB org] \
SMQ_INFLUXDB_TOKEN=[InfluxDB token] \
SMQ_INFLUXDB_DBURL=[InfluxDB database url] \
SMQ_INFLUXDB_USER_AGENT=[InfluxDB user agent] \
SMQ_INFLUXDB_TIMEOUT=[InfluxDB timeout] \
SMQ_INFLUXDB_INSECURE_SKIP_VERIFY=[InfluxDB insecure skip verify] \
SMQ_MESSAGE_BROKER_URL=[Message broker instance URL] \
SMQ_JAEGER_URL=[Jaeger server URL] \
SMQ_SEND_TELEMETRY=[Send telemetry to supermq call home server] \
SMQ_INFLUX_WRITER_INSTANCE_ID=[Influx writer instance ID] \
$GOBIN/supermq-contrib-influxdb
```

### Using docker-compose

This service can be deployed using docker containers.
Docker compose file is available in `<project_root>/docker/addons/influxdb-writer/docker-compose.yml`. Besides database
and writer service, it contains InfluxData Web Admin Interface which can be used for database
exploration and data visualization and analytics. In order to run SupeMQ InfluxDB writer, execute the following command:

```bash
docker compose -f docker/addons/influxdb-writer/docker-compose.yml up -d
```

And, to use the default .env file, execute the following command:

```bash
docker compose -f docker/addons/influxdb-writer/docker-compose.yml up --env-file docker/.env -d
```

_Please note that you need to start core services before the additional ones._

## Usage

Starting service will start consuming normalized messages in SenML format.

Official docs can be found [here](https://docs.supermq.abstractmachines.fr).
