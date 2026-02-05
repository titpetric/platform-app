# Pulse

A client that listens to keypress events, counts them and updates the
count to the server every 5 minutes. You can host your own, or:

- Join [pulse.incubator.to](https://pulse.incubator.to) (WIP)

> TODO: show off a cool graph or two, this is barely an ingest for now.

## Docker

The suggested way how to run the pulse client and server is via docker.
No other system service config is provided at this time.

### Running the server

Assuming you want to run the server locally:

```bash
docker compose up -d pulse
```

The server does nothing on its own, and even exposes no ports. It
defines some labels for "production" (stack specific) which can
safely be deleted for your own use.

### Client authentication

I'd love just a pre-shared secret, but the design of the pulse server allows
for multi tenancy. While you could navigate the web to register, you can also
do that via the following commands.

To create a user:

```bash
docker compose run --rm pulse-client register --server http://pulse:8080
```

The registration already gives a JWT access token. The token gets
rotated on a daily basis and is stored in `$HOME/.config/pulse/token.json`.
The container mounts the location as configured in `compose.yml`.

If you already have a user:

```bash
docker compose run --rm pulse-client login --server http://pulse:8080
```

### Starting the client

After authenticating, to start the client with default options, just run:

```bash
docker compose up -d pulse-client
```

This will send your recorded keystrokes to your local pulse server instance.

### Sending data elsewhere

If you want to send your keystrokes and monitor keystrokes on the public
index, follow the same register/login steps with a custom `--server`.

## Local installation

You can install pulse using `go`:

```bash
go install github.com/titpetric/platform-app/pulse/cmd/pulse@main
```

Run `pulse --help` to discover usage.
