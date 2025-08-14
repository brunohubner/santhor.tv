FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/webserver .

FROM busybox:latest
WORKDIR /app
COPY --from=builder /app/bin/webserver .
EXPOSE ${PORT}
ENTRYPOINT ["./webserver"]