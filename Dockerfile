# CR3 Client Dashboard
# Package: gitlab.hexacode.org/cr3/client-dashboard
# Maintainer: Azzis Arswendo <azzis@hexacode.org>
# 
# Copyright (C) 2024 Hexacode Teknologi Indonesia
# All Rights Reserved


# Build frontend #
FROM oven/bun:latest AS frontend

WORKDIR /app

RUN apt update && apt install -y git python3

COPY frontend /app
COPY config.json /config.json

RUN bun install

ENV NODE_ENV=production

RUN bun run build

RUN mv template/public public


# Build backend #
FROM golang:1.19-alpine AS backend

WORKDIR /app

RUN apk add -U tzdata
RUN apk --update add ca-certificates alpine-sdk


COPY backend /app

RUN go install gitlab.hexacode.org/go-libs/microservice
RUN cp /go/bin/microservice /app
RUN go build -buildmode=plugin -o plugin.so


# Build final image #
FROM alpine:latest

WORKDIR /app

COPY --from=frontend /app/template /app/template
COPY --from=frontend /app/public /app/public

COPY --from=backend /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend /etc/passwd /etc/passwd
COPY --from=backend /etc/group /etc/group

COPY --from=backend /app/plugin.so /app/plugin.so
COPY --from=backend /app/microservice /microservice

# ENV GIN_MODE=release
ENV STRAPI_URL=""
ENV STRAPI_TOKEN=""

EXPOSE 8080

# Copy all the configurations and assets
COPY config.json /app/config.json

CMD [ "/microservice" ]