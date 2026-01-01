# AI Bridges 🚀

**AI Bridges** is a high-performance WebAI-to-API service built in Go. It allows you to convert web-based AI services (like Google Gemini) into standardized REST APIs, including an OpenAI-compatible interface.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://github.com/ntthanh2603/ai-bridges/pkgs/container/ai-bridges)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/ntthanh2603/ai-bridges/blob/main/LICENSE)

---

## ✨ Features

- 🌉 **Service Bridge**: Seamlessly connect web-based AI to your applications.
- 🤖 **Gemini Support**: Full support for Google Gemini (pro) using session cookies.
- 🔄 **Auto Cookie Rotation**: Automatically manages and refreshes session tokens (`__Secure-1PSIDTS`) to keep the connection alive.
- 💾 **Session Persistence**: Caches refreshed cookies locally (`.cookies/`) to survive restarts.
- 🔌 **OpenAI Compatible**: Provides a `/v1/chat/completions` endpoint that mimics OpenAI's API.
- � **Native API Support**: Exposes endpoints compatible with Google's official Gemini `v1beta` REST API.
- �🚀 **Built with Fiber**: Ultra-fast and efficient web framework.
- 📝 **Swagger UI**: Interactive API documentation built-in.
- 🐳 **Dockerized**: Ready for containerized deployment with unified configuration.

---

## 🛠️ Technology Stack

- **Language**: Go (v1.24+)
- **Framework**: [Gofiber/fiber](https://github.com/gofiber/fiber)
- **HTTP Client**: [req/v3](https://github.com/imroc/req/v3)
- **Logging**: Uber-zap
- **Documentation**: Swag / Swagger

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher installed. Or Docker.

### Configuration Priority

The application uses a unified configuration system with the following priority:

1. **Environment Variables** (Highest priority)
2. **`config.yml`**
3. **Defaults** (Lowest priority)

### Environment Variables

| Variable                  | Corresponding YAML Key    | Description                                             |
| ------------------------- | ------------------------- | ------------------------------------------------------- |
| `GEMINI_1PSID`            | `GEMINI_1PSID`            | (Required) Main session cookie                          |
| `GEMINI_1PSIDTS`          | `GEMINI_1PSIDTS`          | (Recommended) Timestamp cookie. Found in browser tools. |
| `GEMINI_1PSIDCC`          | `GEMINI_1PSIDCC`          | (Optional) Context cookie                               |
| `GEMINI_REFRESH_INTERVAL` | `GEMINI_REFRESH_INTERVAL` | Rotation interval in minutes (default: 30)              |
| `PORT`                    | `PORT`                    | Server port (default: 3000)                             |

---

## 🐳 Docker Usage (Quick Start)

The easiest way to get started is to pull the pre-built image.

### 1. Pull the image

```bash
docker pull ghcr.io/ntthanh2603/ai-bridges:latest
```

### 2. Run container

```bash
docker run -d -p 3000:3000 \
  -e GEMINI_1PSID="your_psid_here" \
  -e GEMINI_1PSIDTS="your_psidts_here" \
  -e GEMINI_REFRESH_INTERVAL=30 \
  -v $(pwd)/cookies:/app/.cookies \
  --name ai-bridges \
  --restart unless-stopped \
  ghcr.io/ntthanh2603/ai-bridges:latest
```

---

## 🛠️ Building from Source

If you want to modify the code or run it locally without Docker.

### 1. Clone the repository

```bash
git clone https://github.com/ntthanh2603/ai-bridges.git
cd ai-bridges
```

### 2. Configure

Copy the example config and add your cookies:

```bash
cp config.example.yml config.yml
# Edit config.yml with your GEMINI_1PSID and GEMINI_1PSIDTS
```

### 3. Run the server

```bash
go run cmd/server/main.go
```

---

## 🧪 Quick Testing

Once the server is running, you can test the connection using either the OpenAI-compatible endpoint or the Gemini native endpoint.

### 1. OpenAI Compatible (Preferred)

Compatible with most AI clients (SDKs, LangChain, etc.).

```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini-pro",
    "messages": [{"role": "user", "content": "Hello, who are you?"}]
  }'
```

### 2. Gemini Native (v1beta)

Compatible with Google's Request/Response format.

```bash
curl -X POST \
  'http://localhost:3000/v1beta/models/gemini-pro:generateContent' \
  -H 'Content-Type: application/json' \
  -d '{
    "contents": [{
      "parts": [{"text": "Hello, explain how you work."}]
    }]
  }'
```

---

## 📘 API Documentation

Visit `http://localhost:3000/swagger/` for the full interactive API documentation.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/ntthanh2603/ai-bridges/blob/main/LICENSE) file for details.
