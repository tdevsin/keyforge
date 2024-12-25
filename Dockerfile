# Stage 1: Build the Go binary
FROM golang:1.23.4-alpine3.21 AS build

COPY . /app

WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/keyforge main.go

# Stage 2: Minimal image with the binary
FROM alpine:3.21.0

WORKDIR /app
COPY --from=build /bin/keyforge /app/keyforge

RUN chmod +x /app/keyforge

EXPOSE 8080

# Start command is hardcoded temporarily. This will be changed to a more dynamic approach in the future.
ENTRYPOINT ["./keyforge", "start"]
