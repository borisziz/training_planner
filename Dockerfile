# syntax=docker/dockerfile:1

FROM golang:1.23.0-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/training-planner \
    ./cmd/service

FROM alpine:3.20

WORKDIR /app

COPY --from=build /out/training-planner /app/training-planner

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/app/training-planner"]
