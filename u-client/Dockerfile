# Builder

FROM golang:1.21-alpine AS builder
WORKDIR /usr/local/src

RUN apk --no-cache add bash make git curl gcc musl-dev

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./

CMD go mod download

COPY . ./
RUN go build -o ./bin/u-client cmd/u-client/main.go

# U-client

FROM alpine:latest AS client
COPY --from=builder /usr/local/src/bin/u-client /
COPY config.yaml ./

CMD ["/u-client"]