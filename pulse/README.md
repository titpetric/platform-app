# Pulse

A key counting docker server/client. The server will ingest key counts
from authenticated clients. Clients allow registration and login against
the server. Multiple clients can be added to a single user.

## Configuring the client

Create a `compose.yml` file as such:

```yaml
services:
  pulse-client:
    image: titpetric/pulse
    user: root
    privileged: true
    volumes:
      - $PWD/data:/app
      - /dev/input:/dev/input:ro
    command:
      - "record"
      - "--name"
      - "REPLACE_WITH_DEVICE_NAME"
      - "--server"
      - "http://pulse.incubator.to"
```

Note that reading from keyboard uses:

- `privileged: true` - ability to read from devices,
- `user: root` - the filesystem permissions for dev/input


## Starting the client

A client first needs to authenticate to a server, and store it's
configuration. First, create a `data/` folder and give it write
permissions. This is used to store authentication details for the pulse
client. It stores no other data.

```bash
mkdir -p data
chmod a+rwx data
```

By now you should have:

```text
data/
compose.yml
```

First, create an user on the pulse server:

```bash
docker compose run --rm pulse-client register --server http://pulse.incubator.to
```

And to start sending data start the pulse client:

```bash
docker compose up -d
```

## Running your own server

You can self host your own pulse server.

See [./compose.yml](./compose.yml) and [./compose.client.yml](./compose.client.yml) for server configuration.

Authenticating against the local pulse server:

```bash
docker compose run --rm pulse-client register --server http://pulse:8080
```
