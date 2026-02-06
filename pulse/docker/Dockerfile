FROM alpine:latest

ARG TARGETARCH

RUN apk update && apk add --no-cache coreutils
RUN which uname
RUN id nobody 2>/dev/null || adduser -D -H -u 65534 nobody

COPY ./bin/pulse-linux-${TARGETARCH} /usr/local/bin/pulse

ENV HOME=/app
VOLUME /app
WORKDIR /app

EXPOSE 8080

ENV PLATFORM_DB_USER=sqlite://file:user.db
ENV PLATFORM_DB_PULSE=sqlite://file:pulse.db

USER nobody
ENTRYPOINT ["/usr/local/bin/pulse"]
CMD ["server"]