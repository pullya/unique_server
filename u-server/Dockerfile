# Builder

FROM golang:1.21-alpine AS builder
WORKDIR /usr/local/src

RUN apk --no-cache add bash make git curl gcc musl-dev

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./
RUN go build -o ./bin/u-server cmd/u-server/main.go

# U Server

FROM alpine:latest AS server
COPY --from=builder /usr/local/src/bin/u-server /
COPY config.yaml ./

CMD ["/u-server"]

EXPOSE 8081