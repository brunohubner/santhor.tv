# Build Stage
FROM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -v -o /app/bin/webserver .

RUN apk --no-cache add ca-certificates tzdata


# Final Stage
FROM busybox:latest

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/bin/webserver .

EXPOSE ${PORT}

ENTRYPOINT ["./webserver"]
