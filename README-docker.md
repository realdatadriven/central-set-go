# Central-Set-Go Docker

## Build
```bash
docker build -t central-set-go:latest .
```

## Run

### Start Application
```bash
docker run -d \
  --name central-set-go \
  -p 4444:8080 \
  -v ./.env:/app/.env:ro \
  -v ./database:/app/database \
  central-set-go:latest
```

```bash
docker run -v ./.env:/app/.env:ro -v ./database:/app/database central-set-go:latest --init
```

```bash
docker run -p 8080:4444 -v ./.env:/app/.env:ro -v ./database:/app/database central-set-go:latest
```

### Initialize Database
```bash
docker run --rm \
  -v ~/.env:/app/.env:ro \
  -v ./database:/app/database \
  central-set-go:latest --init
```

### Initialize with Custom Database Name
```bash
docker run --rm \
  -v ~/.env:/app/.env:ro \
  -v ./database:/app/database \
  central-set-go:latest --init --dbname MY_DATABASE
```

## Access
Application runs on `http://localhost:8080`

## Environment
Create `~/.env` file with your configuration or use environment variables directly with `-e` flags.

