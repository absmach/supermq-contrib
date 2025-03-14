# Twins

**Twins** service is used for  creating, retrieving, updating, and deleting digital twins.
A **digital twin**  is an abstract and semantic representation of a real world data system consisting of data producers and consumers.
It stores the sequence of attribute based definitions of a system and refers to a time series of definition based states that store the system historical data.
Digital twin is usually less detailed and can be a digital replica of a real world system such as the industrial machine. It is used to create and store information about system's state at any given moment, to compare system state over a given period of time - so-called diffs or deltas - as well as to control agents composing the system.

## Overview

The Twins Service is built on top of the SupeMQ platform and interacts with its core components such as users, clients, and channels. It listens to the message broker, intercepts relevant messages, and updates digital twin states accordingly. Each twin consists of:

- **Metadata**: owner, ID, name, timestamps, revision number
- **Definitions**: semantic representation of system components as attributes
- **States**: time-series history of system state

SupeMQ twins anatomy is of the following format:

```go
// Twin is a SupeMQ data system representation. Each twin is owned
// by a single user, and is assigned with the unique identifier.
type Twin struct {
 Owner       string
 ID          string
 Name        string
 Created     time.Time
 Updated     time.Time
 Revision    int
 Definitions []Definition
 Metadata    Metadata
}
```

Twin states are persisted in the separate collection of the same database. Currently, twins service uses the MongoDB. InfluxDB support for twins and states persistence is on the roadmap.

When we define our digital twin, its JSON representation might look like this:

```json
{
  "owner": "john.doe@email.net",
  "id": "a838e608-1c1b-4fea-9c34-def877473a89",
  "name": "grinding machine 2",
  "revision": 2,
  "created": "2020-05-05T08:41:39.142Z",
  "updated": "2020-05-05T08:49:12.638Z",
  "definitions": [
    {
      "id": 0,
      "created": "2020-05-05T08:41:39.142Z",
      "attributes": [],
      "delta": 1000000
    },
    {
      "id": 1,
      "created": "2020-05-05T08:46:23.207Z",
      "attributes": [
        {
          "name": "engine temperature",
          "channel": "7ef6c61c-f514-402f-af4b-2401b588bfec",
          "subtopic": "engine",
          "persist_state": true
        },
        {
          "name": "chassis temperature",
          "channel": "7ef6c61c-f514-402f-af4b-2401b588bfec",
          "subtopic": "chassis",
          "persist_state": true
        },
        {
          "name": "rotations per sec",
          "channel": "a254032a-8bb6-4973-a2a1-dbf80f181a86",
          "subtopic": "",
          "persist_state": false
        }
      ],
      "delta": 1000000
    },
    {
      "id": 2,
      "created": "2020-05-05T08:49:12.638Z",
      "attributes": [
        {
          "name": "engine temperature",
          "channel": "7ef6c61c-f514-402f-af4b-2401b588bfec",
          "subtopic": "engine",
          "persist_state": true
        },
        {
          "name": "chassis temperature",
          "channel": "7ef6c61c-f514-402f-af4b-2401b588bfec",
          "subtopic": "chassis",
          "persist_state": true
        },
        {
          "name": "rotations per sec",
          "channel": "a254032a-8bb6-4973-a2a1-dbf80f181a86",
          "subtopic": "",
          "persist_state": false
        },
        {
          "name": "precision",
          "channel": "aed0fbca-0d1d-4b07-834c-c62f31526569",
          "subtopic": "",
          "persist_state": true
        }
      ],
      "delta": 1000000
    }
  ]
}
```

In the case of the upper twin, we begin with an empty definition, the one with the `id` **0** - we could have provided the definition immediately - and over the course of time, we add two more definitions, so the total number of revisions is **2** (revision index is zero-based). We decide not to persist the number of rotation per second in our digital twin state. We define it, though, because the definition and its attributes are used not only to define states of a complex data agent system, but also to define the semantic structure of the system. `delta` is the number of nanoseconds used to determine whether the received attribute value should trigger the generation of the new state or the same state should be updated. The reason for this is to enable state sampling over the regular intervals of time. Discarded values are written to the database of choice by SupeMQ [writers][writer], so you can always retrieve intermediate values if need be.

**states** are created according to the twin's current definition. A state stores twin's ID - every state belongs to a single twin -, its own ID, twin's definition number, creation date and the actual payload. **Payload** is a set of key-value pairs where a key corresponds to the attribute name and a value is the actual value of the attribute. All [SenML value types][senml] are supported.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                   | Description                                                         | Default                          |
| -------------------------- | ------------------------------------------------------------------- | -------------------------------- |
| SMQ_TWINS_LOG_LEVEL         | Log level for twin service (debug, info, warn, error)               | info                             |
| SMQ_TWINS_HTTP_PORT         | Twins service HTTP port                                             | 9018                             |
| SMQ_TWINS_SERVER_CERT       | Path to server certificate in PEM format                            |                                  |
| SMQ_TWINS_SERVER_KEY        | Path to server key in PEM format                                    |                                  |
| SMQ_JAEGER_URL              | Jaeger server URL                                                   | <http://jaeger:14268/api/traces> |
| SMQ_TWINS_DB                | Database name                                                       | supermq                       |
| SMQ_TWINS_DB_HOST           | Database host address                                               | localhost                        |
| SMQ_TWINS_DB_PORT           | Database host port                                                  | 27017                            |
| SMQ_CLIENTS_STANDALONE_ID    | User ID for standalone mode (no gRPC communication with users)      |                                  |
| SMQ_CLIENTS_STANDALONE_TOKEN | User token for standalone mode that should be passed in auth header |                                  |
| SMQ_TWINS_CLIENT_TLS        | Flag that indicates if TLS should be turned on                      | false                            |
| SMQ_TWINS_CA_CERTS          | Path to trusted CAs in PEM format                                   |                                  |
| SMQ_TWINS_CHANNEL_ID        | Message broker notifications channel ID                             |                                  |
| SMQ_MESSAGE_BROKER_URL      | SupeMQ Message broker URL                                       | <nats://localhost:4222>          |
| SMQ_AUTH_GRPC_URL           | Auth service gRPC URL                                               | <localhost:7001>                 |
| SMQ_AUTH_GRPC_TIMEOUT       | Auth service gRPC request timeout in seconds                        | 1s                               |
| SMQ_TWINS_CACHE_URL         | Cache database URL                                                  | <redis://localhost:6379/0>       |
| SMQ_SEND_TELEMETRY          | Send telemetry to supermq call home server                       | true                             |

## Deployment

The service itself is distributed as Docker container. Check the [`twins`](https://github.com/absmach/supermq-contrib/blob/main/docker/addons/twins/docker-compose.yml#L35-L58) service section in
docker-compose file to see how service is deployed.

To start the service outside of the container, execute the following shell
script:

```bash
# download the latest version of the service
go get github.com/absmach/supermq-contrib

cd $GOPATH/src/github.com/absmach/supermq-contrib

# compile the twins
make twins

# copy binary to bin
make install

# set the environment variables and run the service
SMQ_TWINS_LOG_LEVEL=[Twins log level] \
SMQ_TWINS_HTTP_PORT=[Service HTTP port] \
SMQ_TWINS_SERVER_CERT=[String path to server cert in pem format] \
SMQ_TWINS_SERVER_KEY=[String path to server key in pem format] \
SMQ_JAEGER_URL=[Jaeger server URL] \
SMQ_TWINS_DB=[Database name] \
SMQ_TWINS_DB_HOST=[Database host address] \
SMQ_TWINS_DB_PORT=[Database host port] \
SMQ_CLIENTS_STANDALONE_EMAIL=[User email for standalone mode (no gRPC communication with auth)] \
SMQ_CLIENTS_STANDALONE_TOKEN=[User token for standalone mode that should be passed in auth header] \
SMQ_TWINS_CLIENT_TLS=[Flag that indicates if TLS should be turned on] \
SMQ_TWINS_CA_CERTS=[Path to trusted CAs in PEM format] \
SMQ_TWINS_CHANNEL_ID=[Message broker notifications channel ID] \
SMQ_MESSAGE_BROKER_URL=[SupeMQ Message broker URL] \
SMQ_AUTH_GRPC_URL=[Auth service gRPC URL] \
SMQ_AUTH_GRPC_TIMEOUT=[Auth service gRPC request timeout in seconds] \
SMQ_TWINS_CACHE_URL=[Cache database URL] \
$GOBIN/supermq-contrib-twins
```

## API Usage

### Starting twins service

The twins service publishes notifications on a Message broker subject of the format
`channels.<SMQ_TWINS_CHANNEL_ID>.messages.<twinID>.<crudOp>`, where `crudOp`
stands for the crud operation done on twin - create, update, delete or
retrieve - or state - save state. In order to use twin service notifications,
one must inform it - via environment variables - about the SupeMQ channel used
for notification publishing. You must use an already existing channel, since you
cannot know in advance or set the channel ID (SupeMQ does it automatically).

To set the environment variable, please go to `.env` file and set the following
variable:

```bash
SMQ_TWINS_CHANNEL_ID=
```

### Create a Twin

Create and update requests use JSON body to initialize and modify, respectively, twin. You can omit every piece of data - every key-value pair - from the JSON. However, you must send at least an empty JSON body.

Create request uses POST HTTP method to create twin:

```bash
curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer <user_token>"   http://localhost:9018/twins -d '{ "name": "twin_name", "definition": { "attributes": [ { "name": "temperature", "channel": "3b57b952-318e-47b5-b0d7-a14f61ecd03b", "subtopic": "temperature", "persist_state": true } ], "delta": 1 } }'
```

If you do not suply the definition, the empty definition of the form

```json
{
  "id": 0,
  "created": "2020-05-05T08:41:39.142Z",
  "attributes": [],
  "delta": 1000000
}
```

will be created.

### Retrieve a Twin

To view a specific twin:

```bash
curl -s -X GET -H "Authorization: Bearer <user_token>"   http://localhost:9018/twins/ <twin_id>
```

### List Twins

```bash
curl -s -X GET -H "Authorization: Bearer <user_token>"   http://localhost:9018/twins?offset=10&limit=20
```

List requests accept `limit` and `offset` query parameters. By default, i.e. without these parameters, list requests fetches only first ten twins (or less, if there are less then ten twins).

### Update a Twin

```bash
curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer <user_token>"   http://localhost:9018/twins/<twin_id> -d '<twin_data>'
```

### Delete a Twin

```bash
curl -s -X DELETE -H "Authorization: Bearer <user_token>"   http://localhost:9018/twins/<twin_id>
```

### Fetch Twin States

```bash
curl -s -X GET -H "Authorization: Bearer <user_token>"   http://localhost:9018/states/<twin_id>?offset=10&limit=20
```

## Notifications

Twins service publishes notifications to a SupeMQ message broker channel.

In order to pick up this notification, you have to create a SupeMQ channel before you start the twins service and inform the twins service about the channel by means of the environment variable, like this:

```bash
export SMQ_TWINS_CHANNEL_ID=f6894dfe-a7c9-4eef-a614-637ebeea5b4c
```

The twins service will use this channel to publish notifications related to twins creation, update, retrieval and deletion. It will also publish notifications related to state saving into the database.

All notifications will be published on the following message broker subject:

```txt
channels.<SMQ_TWINS_CHANNEL_ID>.<optional_subtopic>
```

where `<optional_subtopic>` is one of the following:

- `create.success` - on successful twin creation,
- `create.failure` - on twin creation failure,
- `update.success` - on successful twin update,
- `update.failure` - on twin update failure,
- `get.success` - on successful twin retrieval,
- `get.failure` - on twin retrieval failure,
- `remove.success` - on successful twin deletion,
- `remove.failure` - on twin deletion failure,
- `save.success` - on successful state save
- `save.failure` - on state save failure.

## Authentication & Authorization

Each twin belongs to a SupeMQ user (a person or an organization).
API calls require an authentication token (`Bearer <user_token>` in the header).

## Additional Resources

with the corresponding values of the desired channel. If you are running
supermq natively, than do the same client in the corresponding console
environment.

For more details, visit the [API documentation](https://docs.api.supermq.abstractmachines.fr/?urls.primaryName=twins-openapi.yml).

[writer]: ./storage.md
[senml]: https://tools.ietf.org/html/rfc8428#section-4.3
