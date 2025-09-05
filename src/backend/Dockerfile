
# ---- Builder ----
    FROM golang:1.24.4-bookworm AS builder
    ENV GOTOOLCHAIN=auto
    
    WORKDIR /app
    
    # Copy module files first (cache)
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the project
    COPY . .
    
    # Static build (binary at /bin/netsecure-iq)
    RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/netsecure-iq .
    
    # ---- Runtime ----
    FROM debian:bookworm-slim
    
    # wg and docker CLI for the `wg genkey` and `docker exec` calls performed by the backend
    RUN apt-get update && apt-get install -y --no-install-recommends \
          wireguard-tools docker.io ca-certificates tzdata && \
        rm -rf /var/lib/apt/lists/*
    
    WORKDIR /app
    
    # Binary
    COPY --from=builder /bin/netsecure-iq /app/backend
    
    # Port exposed by main.go
    EXPOSE 8000
    
    # IMPORTANT: if the backend calls `docker exec`, mount the host socket at run:
    #   -v /var/run/docker.sock:/var/run/docker.sock
    CMD ["/app/backend"]
    