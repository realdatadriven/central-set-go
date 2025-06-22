# ============================================
# 🛠️ Stage 1: Build central-set-go from Source
# ============================================
FROM golang:1.24 as builder

# Set working directory inside the container
WORKDIR /app

# Install system dependencies required for building
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    unixodbc \
    unixodbc-dev \
    && rm -rf /var/lib/apt/lists/*

# Enable CGO for ODBC support
ENV CGO_ENABLED=1

# Clone the central-set-go repository
RUN git clone --depth=1 https://github.com/realdatadriven/central-set-go.git .

# Build the central-set-go binary
RUN go build -o central-set-go ./cmd/api/main.go

# ============================================
# 🚀 Stage 2: Create Minimal Runtime Image
# ============================================
FROM debian:bookworm-slim

# Install runtime dependencies (unixODBC)
RUN apt-get update && apt-get install -y \
    ca-certificates \
    unixodbc \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the compiled central-set-go binary from the builder stage
COPY --from=builder /app/central-set-go /usr/local/bin/central-set-go

# Ensure the binary is executable
RUN chmod +x /usr/local/bin/central-set-go

# Allow users to mount a config file
VOLUME ["/app/config", "/app/data"]

# Set the entrypoint to pass CLI arguments
ENTRYPOINT ["/usr/local/bin/central-set-go"]

#docker build -t central-set-go:latest .

