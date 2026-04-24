# build stage
FROM golang:1.26 AS builder

WORKDIR /app

# copy only what we need
COPY ./cmd ./cmd
COPY ./gen ./gen
COPY ./internal ./internal
COPY go.mod go.sum Makefile ./

RUN CGO_ENABLED=0 GOOS=linux make build

# final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/pack-optimizer ./pack-optimizer

EXPOSE 8080

CMD ["./pack-optimizer"]