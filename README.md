# Ollama CORS Proxy

A lightweight proxy server that adds CORS headers to Ollama API responses.

## Prerequisites
- Docker
- Ollama running locally or accessible URL

### Installing Ollama
1. Install Ollama using curl (macOS/Linux):
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

2. Start the Ollama service:
```bash
ollama serve
```

3. Verify the installation:
```bash
ollama --version
```

The Ollama API will be available at `http://localhost:11434` by default.

## Quick Start

### Using Docker Compose
```yaml
version: '3'
services:
  ollama-cors-proxy:
    build: .
    ports:
      - "11435:11435"
    environment:
      - OLLAMA_URL=http://host.docker.internal:11434
```

Start the service:
```bash
docker-compose up -d
```

> Note: Use `host.docker.internal` to connect to Ollama running on your host machine.

## License
MIT