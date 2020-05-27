# syntax = docker/dockerfile:experimental

FROM bitnami/kubectl:1.17.4 as kubectl

FROM golang:1.14 as dev
WORKDIR /app

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl
# install dev tools
RUN apt-get update && \
  apt-get install postgresql-client -y && \ 
  rm -rf /var/lib/apt/lists/*

RUN go get github.com/pilu/fresh && \
  go get github.com/go-delve/delve/cmd/dlv

# preload dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM dev as build
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -v -o guestbook .

FROM alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true


WORKDIR /app
COPY --from=build /app/guestbook /app/guestbook
COPY --from=build /app/public /app/public
EXPOSE 8080
CMD [ "/app/guestbook" ]

