# Pulse

A client that listens to keypress events, counts them and updates the
count to the server every 5 minutes. You can host your own, or:

- Join [pulse.incubator.to](https://pulse.incubator.to)

TODO: show off a cool graph or two, this is barely an ingest for now.

## Local Client/Server setup

Asuming you want to run the server locally:

```bash
docker compose up -d pulse
```

The server does nothing on it's own, and even exposes no ports. It
defines some labels for "production" (stack specific) which can
safely be deleted for your own use.

## Client authentication

I'd love just a pre-shared secret, but the design of the pulse server allows
for multi tenancy. While you could navigate the web to register, you can also
do that via the following commands.

To create a user:

```bash
docker compose run --rm pulse-client register --server http://pulse:8080
```

The registration already gives a JWT access token. The token gets
rotated on a daily basis and is stored in `$HOME/.config/pulse.json`.
The container mounts the location as configured in `compose.yml`.

If you already have a user:

```bash
docker compose run --rm pulse-client login --server http://pulse:8080
```

## Starting the client

After authenticating, to start the client with default options, just run:

```
docker compose up -d
```

This will send your recorded keystrokes to your local pulse server instance.

## Sending data elsewhere

If you want to send your keystrokes and monitor keystrokes on the public
index, follow the same register/login steps with a custom `--server`.
