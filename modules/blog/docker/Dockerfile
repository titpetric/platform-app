FROM alpine:latest

ENV PLATFORM_DB_BLOG=sqlite://file:cache/blog.db

ADD ./blog /usr/local/bin/blog

ENTRYPOINT ["/usr/local/bin/blog"]

WORKDIR /app
