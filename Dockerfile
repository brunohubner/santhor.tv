FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /webserver .
RUN apk --no-cache add ca-certificates tzdata

FROM busybox:latest
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/webserver .
EXPOSE ${PORT}
ENTRYPOINT ["./webserver"]
