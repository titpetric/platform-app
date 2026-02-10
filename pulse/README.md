# Pulse

A key counting docker server/client.

## Configuring the client

Create a `compose.yml` file as such:

```
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
      - "chronos"
      - "--server"
      - "http://pulse.incubator.to"
      - "--duration"
      - "5m"
```

Not great isolation when you want to read from input devices.

## Starting the client

A client first needs to authenticate to a server, and store it's
configuration. First, create a `data/` folder and give it write
permissions.

```bash
mkdir -p data
chmod a+rwx data
```

You can register (or login) via the following commands:

```bash
docker compose run --rm pulse-client register --server http://pulse.incubator.to
docker compose up -d
```

## Running your own server

See [./compose.yml](./compose.yml) for server configuration.