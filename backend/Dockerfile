FROM golang:1.25.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /src/backend
RUN go mod download

COPY backend/ ./

ARG TARGET_CMD=./cmd/api
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/app ${TARGET_CMD}

FROM alpine:3.22

LABEL org.opencontainers.image.source="https://github.com/crowdfunding/crowdfunding"
LABEL org.opencontainers.image.description="CrowdFund DApp Backend"

RUN apk add --no-cache ca-certificates tzdata curl

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /out/app /app/app
COPY backend/config /app/config
COPY backend/migrations /app/migrations

RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/app"]
