# cfProxyHub

A modern web dashboard for managing Cloudflare accounts, tunnels, and zones through an intuitive interface. Built with Go and featuring a responsive web UI for streamlined Cloudflare infrastructure management.

## ✨ Features

- **🔐 Secure Authentication** - Admin login with session management
- **🌐 Multi-Account Management** - Manage multiple Cloudflare accounts from one dashboard
- **🚇 Tunnel Management** - Create, configure, and monitor Cloudflare Tunnels
- **📍 Zone Management** - DNS zone configuration and management
- **🔗 Public Hostname Configuration** - Easy setup of tunnel public hostnames
- **📊 Real-time Dashboard** - Monitor your Cloudflare infrastructure at a glance
- **🎨 Modern UI** - Beautiful, responsive web interface built with Corona Bootstrap admin template

## 🏗️ Architecture

```
cfProxyHub/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Data models
│   ├── routes/          # Route definitions
│   └── services/        # Business logic
├── pkg/utils/           # Utility functions
├── web/                 # Frontend assets and templates
│   ├── assets/          # CSS, JS, images
│   └── templates/       # HTML templates
└── tests/               # Test files
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24+ installed
- Cloudflare account with API access
- API Token or API Key + Email

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/cfProxyHub.git
   cd cfProxyHub
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   
   Create a `.env` file in the project root:
   ```env
   # Cloudflare Credentials (choose one method)
   CLOUDFLARE_API_TOKEN=your_api_token_here
   # OR
   CLOUDFLARE_API_KEY=your_api_key_here
   CLOUDFLARE_EMAIL=your_cloudflare_email@example.com
   
   # Server Configuration
   PORT=8080
   
   # Admin Authentication
   ADMIN_USERNAME=admin
   ADMIN_PASSWORD=your_secure_password
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the dashboard**
   
   Open your browser and navigate to: `http://localhost:8080`

## 📋 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `CLOUDFLARE_API_TOKEN` | Cloudflare API Token | Yes* | - |
| `CLOUDFLARE_API_KEY` | Cloudflare API Key | Yes* | - |
| `CLOUDFLARE_EMAIL` | Cloudflare account email | Yes* | - |
| `PORT` | Server port | No | `8080` |
| `ADMIN_USERNAME` | Admin username | No | `admin` |
| `ADMIN_PASSWORD` | Admin password | No | `password123` |

*Either `CLOUDFLARE_API_TOKEN` OR both `CLOUDFLARE_API_KEY` and `CLOUDFLARE_EMAIL` are required.

### Cloudflare API Setup

#### Option 1: API Token (Recommended)
1. Go to [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
2. Click "Create Token"
3. Use the "Custom token" template
4. Set permissions:
   - Zone:Zone:Read
   - Zone:DNS:Edit
   - Account:Cloudflare Tunnel:Edit

#### Option 2: Global API Key
1. Go to [Cloudflare API Keys](https://dash.cloudflare.com/profile/api-tokens)
2. Copy your Global API Key
3. Use your Cloudflare email address

## 🛠️ API Endpoints

### Authentication
- `POST /api/auth/login` - Admin login
- `POST /api/auth/logout` - Admin logout

### Cloudflare Management
- `GET /api/cloudflare/accounts` - List all accounts
- `GET /api/cloudflare/accounts/:id/tunnels` - Get tunnels for account
- `POST /api/cloudflare/accounts/:id/tunnels` - Create new tunnel
- `GET /api/cloudflare/accounts/:id/zones` - Get zones for account
- `POST /api/cloudflare/tunnels/:id/hostnames` - Create public hostname

### Web Interface
- `GET /` - Dashboard (requires authentication)
- `GET /login` - Login page
- `GET /tunnels` - Tunnel management page
- `GET /accounts` - Account management page

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

Run specific tests:
```bash
go test ./tests/
```

## 🏗️ Development

### Project Structure

- **`cmd/server/`** - Application entry point and main function
- **`internal/`** - Private application code
  - **`config/`** - Configuration loading and validation
  - **`handlers/`** - HTTP request handlers for different endpoints
  - **`middleware/`** - Authentication and other middleware
  - **`models/`** - Data structures and models
  - **`routes/`** - Route definitions and setup
  - **`services/`** - Business logic and Cloudflare API integration
- **`pkg/`** - Public utilities and helpers
- **`web/`** - Frontend assets (CSS, JS, images, templates)
- **`tests/`** - Test files

### Building for Production

```bash
# Build binary
go build -o cfproxyhub cmd/server/main.go

# Run binary
./cfproxyhub
```

### Docker Deployment

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o cfproxyhub cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/cfproxyhub .
COPY --from=builder /app/web ./web
EXPOSE 8080
CMD ["./cfproxyhub"]
```

## 🔒 Security Considerations

- Always use strong passwords for admin authentication
- Use API tokens with minimal required permissions
- Run behind a reverse proxy (nginx, Cloudflare) in production
- Keep your Cloudflare credentials secure
- Regularly rotate API tokens

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cloudflare Go SDK](https://github.com/cloudflare/cloudflare-go) for Cloudflare API integration
- [Gin Web Framework](https://github.com/gin-gonic/gin) for HTTP routing
- [godotenv](https://github.com/joho/godotenv) for environment variable management
- [Corona Admin Dashboard](https://themewagon.com/themes/corona-free-responsive-bootstrap-4-html-5-admin-dashboard-template/) by ThemeWagon for the beautiful web interface template

## 📞 Support

If you have any questions or need help, please:
1. Check the [Issues](https://github.com/yourusername/cfProxyHub/issues) for existing solutions
2. Create a new issue if your problem isn't covered
3. Provide detailed information about your setup and the issue

---

**Made with ❤️ for simplified Cloudflare management**