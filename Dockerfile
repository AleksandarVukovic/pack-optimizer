# build stage
FROM golang:1.26 AS builder

WORKDIR /app

# copy only what we need
COPY ./cmd ./cmd
COPY ./gen ./gen
COPY ./internal ./internal
COPY go.mod go.sum Makefile ./
COPY ./web/static/index.html ./index.html

RUN CGO_ENABLED=0 GOOS=linux make build

# final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/pack-optimizer ./pack-optimizer
COPY --from=builder /app/index.html ./index.html

ENV PORT=8080

EXPOSE $PORT

CMD ["./pack-optimizer"]