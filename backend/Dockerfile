# Stage 1 - Builder
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/flagr \
    ./cmd/server/

# Stage 2 - Runtime
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/flagr /flagr
COPY --from=builder /app/migrations /migrations

EXPOSE 8080 50051

USER 1001

ENTRYPOINT ["/flagr"]