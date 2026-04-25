# build stage
FROM golang:1.26 AS builder

WORKDIR /app

# copy only what we need
COPY ./cmd ./cmd
COPY ./gen ./gen
COPY ./internal ./internal
COPY ./web/static ./static
COPY go.mod go.sum Makefile ./

RUN CGO_ENABLED=0 GOOS=linux make build

# final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/pack-optimizer ./pack-optimizer

# copy index and docs files
COPY --from=builder /app/static ./
COPY --from=builder /app/gen/http/openapi3.json ./openapi3.json


ENV PORT=8080

EXPOSE $PORT

CMD ["./pack-optimizer"]