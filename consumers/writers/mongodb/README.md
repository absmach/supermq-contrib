# MongoDB writer

MongoDB writer provides message repository implementation for MongoDB.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                         | Description                                                                       | Default                        |
| -------------------------------- | --------------------------------------------------------------------------------- | ------------------------------ |
| SMQ_MONGO_WRITER_LOG_LEVEL        | Log level for MongoDB writer                                                      | info                           |
| SMQ_MONGO_WRITER_CONFIG_PATH      | Config file path with Message broker subjects list, payload type and content-type | /config.toml                   |
| SMQ_MONGO_WRITER_HTTP_HOST        | Service HTTP host                                                                 | localhost                      |
| SMQ_MONGO_WRITER_HTTP_PORT        | Service HTTP port                                                                 | 9010                           |
| SMQ_MONGO_WRITER_HTTP_SERVER_CERT | Service HTTP server certificate path                                              | ""                             |
| SMQ_MONGO_WRITER_HTTP_SERVER_KEY  | Service HTTP server key                                                           | ""                             |
| SMQ_MONGO_NAME                    | Default MongoDB database name                                                     | messages                       |
| SMQ_MONGO_HOST                    | Default MongoDB database host                                                     | localhost                      |
| SMQ_MONGO_PORT                    | Default MongoDB database port                                                     | 27017                          |
| SMQ_MESSAGE_BROKER_URL            | Message broker instance URL                                                       | nats://localhost:4222          |
| SMQ_JAEGER_URL                    | Jaeger server URL                                                                 | http://jaeger:14268/api/traces |
| SMQ_SEND_TELEMETRY                | Send telemetry to supermq call home server                                     | true                           |
| SMQ_MONGO_WRITER_INSTANCE_ID      | MongoDB writer instance ID                                                        | ""                             |

## Deployment

The service itself is distributed as Docker container. Check the [`mongodb-writer`](https://github.com/absmach/supermq-contrib/blob/main/docker/addons/mongodb-writer/docker-compose.yml#L36-L55) service section in docker-compose file to see how service is deployed.

To start the service, execute the following shell script:

```bash
# download the latest version of the service
git clone https://github.com/absmach/supermq-contrib.git

cd supermq-contrib

# compile the mongodb writer
make mongodb-writer

# copy binary to bin
make install

# Set the environment variables and run the service
SMQ_MONGO_WRITER_LOG_LEVEL=[MongoDB writer log level] \
SMQ_MONGO_WRITER_CONFIG_PATH=[Configuration file path with Message broker subjects list] \
SMQ_MONGO_WRITER_HTTP_HOST=[Service HTTP host] \
SMQ_MONGO_WRITER_HTTP_PORT=[Service HTTP port] \
SMQ_MONGO_WRITER_HTTP_SERVER_CERT=[Service HTTP server certificate] \
SMQ_MONGO_WRITER_HTTP_SERVER_KEY=[Service HTTP server key] \
SMQ_MONGO_NAME=[MongoDB database name] \
SMQ_MONGO_HOST=[MongoDB database host] \
SMQ_MONGO_PORT=[MongoDB database port] \
SMQ_MESSAGE_BROKER_URL=[Message broker instance URL] \
SMQ_JAEGER_URL=[Jaeger server URL] \
SMQ_SEND_TELEMETRY=[Send telemetry to supermq call home server] \
SMQ_MONGO_WRITER_INSTANCE_ID=[MongoDB writer instance ID] \

$GOBIN/supermq-contrib-mongodb-writer
```

## Usage

Starting service will start consuming normalized messages in SenML format.
