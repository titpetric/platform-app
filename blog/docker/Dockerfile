FROM alpine:latest

ARG NAME
ARG TARGETOS
ARG TARGETARCH

RUN apk update && apk add --no-cache coreutils git
RUN which uname
RUN id nobody 2>/dev/null || adduser -D -H -u 65534 nobody

COPY ./bin/${NAME}-${TARGETOS}-${TARGETARCH} /usr/local/bin/${NAME}

ENV HOME=/app
WORKDIR /app
RUN chown nobody:nobody /app && touch /app/.keep
VOLUME /app

EXPOSE 8080

# This section is custom for app type
ENV PLATFORM_DB_USER=sqlite://file:user.db
ENV PLATFORM_DB_BLOG=sqlite://file:blog.db

# Create a small wrapper with the binary path hardcoded
RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo "exec /usr/local/bin/$NAME \"\$@\"" >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

USER nobody
ENTRYPOINT ["/entrypoint.sh"]
CMD ["server"]
