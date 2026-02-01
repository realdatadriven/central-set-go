# ============================================
# 🛠️ Stage 1: Build central-set-go from Source
# ============================================
FROM golang:1.25 as builder

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

# Build the central-set-go binary with static linking to avoid GLIBC issues
RUN go build -tags="duckdb_arrow" -o central-set-go ./cmd/api
RUN rm ./database/ADMIN.db

# ============================================
# 🚀 Stage 2: Create Minimal Runtime Image
# ============================================
# Use Ubuntu 24.04 which has GLIBC 2.39, compatible with Go 1.24
FROM ubuntu:24.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    unixodbc \
    && rm -rf /var/lib/apt/lists/*


# Set working directory
WORKDIR /app

# Copy the compiled central-set-go binary from the builder stage
COPY --from=builder /app/central-set-go /usr/local/bin/central-set-go

# Copy static folder from the builder stage (if it exists)
COPY --from=builder /app/static /app/static
RUN mkdir -p /app/static/uploads
COPY --from=builder /app/database /app/database
COPY --from=builder /app/database /app/database.defaults
COPY --from=builder /app/locales /app/locales

# Ensure the binary is executable
RUN chmod +x /usr/local/bin/central-set-go

# Create directory for database
#RUN mkdir -p /app/database

# Define volume for database only
VOLUME ["/app/database"]
#, "/app/static/uploads"]

# Expose common ports (adjust as needed for your application)
EXPOSE 4444

# Create a flexible entrypoint script inline
RUN echo '#!/bin/bash\n\
# Set default environment variables if not provided\n\
export APP_ENV=${APP_ENV:-production}\n\
export DATABASE_PATH=${DATABASE_PATH:-/app/database}\n\
\n\
# Check if .env file exists and source it\n\
if [ -f "/app/.env" ]; then\n\
    echo "Loading environment variables from /app/.env"\n\
    set -a\n\
    source /app/.env\n\
    set +a\n\
fi\n\
# Bootstrap database directory if empty or missing files
if [ -d "/app/database.defaults" ]; then\n\
  for f in /app/database.defaults/*; do\n\
    name=$(basename "$f")\n\
    if [ ! -e "/app/database/$name" ]; then\n\
      echo "Bootstrapping $name"\n\
      cp -a "$f" "/app/database/$name"\n\
    fi\n\
  done\n\
fi\n\
\n\
# Ensure database directory exists\n\
mkdir -p "$DATABASE_PATH"\n\
\n\
# Handle different command scenarios\n\
case "$1" in\n\
    --help|-h)\n\
        echo "Usage: docker run [OPTIONS] central-set-go:latest [COMMAND] [ARGS]"\n\
        echo "Commands:"\n\
        echo "  --init                    Initialize the application"\n\
        echo "  --init --dbname NAME      Initialize with specific database name"\n\
        echo "  --help                    Show this help message"\n\
        echo "  [no args]                 Start the application normally"\n\
        exit 0\n\
        ;;\n\
    "")\n\
        echo "Starting central-set-go in normal mode..."\n\
        exec /usr/local/bin/central-set-go\n\
        ;;\n\
    *)\n\
        echo "Executing: /usr/local/bin/central-set-go $@"\n\
        exec /usr/local/bin/central-set-go "$@"\n\
        ;;\n\
esac' > /entrypoint.sh && chmod +x /entrypoint.sh

# Use the entrypoint script to handle initialization and regular commands
ENTRYPOINT ["/entrypoint.sh"]

# Default command (can be overridden)
CMD [""]
# ============================================
# 📝 Usage Instructions
#docker build --no-cache -t central-set-go:latest .
#docker run -v ./database:/app/database central-set-go:latest --init
#docker run -p 8080:4444 -v ./.env:/app/.env:ro -v ./database:/app/database central-set-go:latest
#podman tag central-set-go:latest docker.io/realdatadriven/central-set-go:latest
#podman tag central-set-go:latest docker.io/realdatadriven/central-set-go:v1.1.8
#podman login docker.io
#podman push docker.io/realdatadriven/central-set-go:latest
#podman push docker.io/realdatadriven/central-set-go:v1.1.8
#docker exec -it c78f3f267461 bash

