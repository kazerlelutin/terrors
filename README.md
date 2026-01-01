# 🎭 Terrors - Error Monitoring Service

> _"The call is coming from inside the house..."_

A horror-themed JavaScript error monitoring service written in Go. When your application crashes, we'll be there to capture the screams... I mean, the errors.

_"In space no one can hear you scream, but in our logs, every error is preserved forever."_

## 🏚️ Project Structure

_"They're here..." - The files are watching you_

```
terrors/
├── cmd/
│   ├── migrate/
│   │   └── main.go          # Migration tool for database schema
│   └── server/
│       └── main.go          # The main entrance to our haunted mansion
├── internal/
│   ├── api/
│   │   ├── handlers.go      # The handlers that never sleep
│   │   └── middleware.go    # Authentication middleware
│   ├── database/
│   │   └── db.go           # The basement where errors are stored
│   ├── models/
│   │   ├── app.go          # Application models
│   │   ├── error.go        # The ghosts of your past mistakes
│   │   ├── webhook.go      # Webhook models
│   │   └── command.go      # CQRS command models
│   └── services/
│       ├── app.go          # Application management service
│       ├── error.go       # Error management service
│       └── webhook.go     # Webhook service
├── migrations/
│   ├── 001_initial_schema.sql
│   └── 002_add_origins_and_webhooks.sql
├── static/
│   └── terrors.js          # The script that haunts your browser
├── go.mod                  # The spell book of dependencies
└── README.md              # This very document you're reading
```

## 🔪 Installation

_"Come play with us, forever and ever..."_

1. **Clone the repository**

```bash
git clone <repository>
cd terrors
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Install Air (hot reload)** (optional)

```bash
go install github.com/air-verse/air@latest
```

4. **Configure environment variables**

Copy the example file:

```bash
cp env.example .env
```

Then modify `.env` according to your configuration (see [Environment Variables](#-environment-variables)).

5. **Run database migrations**

```bash
go run cmd/migrate/main.go
```

6. **Start the server**

**Development mode (with hot reload):**

```bash
air
```

**Production mode:**

```bash
go run cmd/server/main.go
```

The server will be accessible at `http://localhost:3000`

_"The Overlook Hotel has been waiting for you..."_

## 🎬 API Endpoints

_"What's your favorite scary movie?"_

### Public Endpoints

- `GET /` - Welcome to the Overlook Hotel
- `POST /sadako` - Frontend error capture (JavaScript errors from browser)
- `POST /jason` - Backend error capture (server-side errors)
- `GET /cdn/terrors.js` - The script that haunts your browser
- `GET /overlook` - Dashboard CQRS - List all available commands

### Protected Endpoints (require `ADMIN_TOKEN`)

#### Applications Management

- `GET /api/apps` - List all applications
- `POST /api/apps` - Create a new application
- `GET /api/apps/{appId}` - Get application details
- `PATCH /api/apps/{appId}` - Update application (name, description, origins, status)
- `DELETE /api/apps/{appId}` - Deactivate application (soft delete)

#### Errors Management

- `GET /api/errors?appId={appId}&status={status}&limit={limit}` - List errors
- `GET /api/errors/{id}` - Get error details
- `PATCH /api/errors/{id}` - Update error status (`new`, `treated`, `deleted`)
- `GET /api/errors/stats?appId={appId}` - Get error statistics

#### Webhooks Management

- `GET /api/webhooks?appId={appId}` - List webhooks for an application
- `POST /api/webhooks?appId={appId}` - Create a webhook (Discord or GitHub)
- `DELETE /api/webhooks/{id}` - Delete a webhook

## 👻 Client Script Usage

_"Don't fall asleep..."_

### Frontend (JavaScript)

```html
<script
  src="http://localhost:3000/cdn/terrors.js"
  data-app-id="app_xxxxxxxx"
></script>
```

_"The script will watch your application and capture every error that dares to appear."_

### Backend (Any Language)

Capturing backend errors is simple - just send a `POST` to `/jason` (different from `/sadako` which is for frontend):

**Simple example (curl):**

```bash
curl -X POST http://localhost:3000/jason \
  -H "Content-Type: application/json" \
  -d '{
    "appId": "app_xxxxxxxx",
    "message": "Database connection failed",
    "stack": "Error: connection timeout\n    at connect (db.js:42:11)",
    "fingerprint": "a1b2c3d4e5f6...",
    "url": "http://localhost:8080/api/users",
    "ts": 1234567890,
    "type": "error"
  }'
```

**Ready-to-use examples:**

- **Go**: See `examples/backend-go/main.go` for a complete client with panic protection
- **Node.js**: See `examples/backend-nodejs/terrors.js` for Express/Fastify middleware
- **Bun**: See `examples/backend-bun/` for TypeScript client with `Bun.serve()` integration
- **Other languages**: See `examples/README.md` for integration patterns

The endpoint accepts the same format as the frontend script, making it easy to capture errors from any backend!

**Important:** The `app_id` must be created first via the API (see [Creating an Application](#creating-an-application)).

## 🔮 Environment Variables

_"The power of Christ compels you..." to configure these properly_

Copy `env.example` to `.env` and configure:

### **Server Configuration**

- `PORT` - Server port (default: 3000)
- `ADMIN_TOKEN` - **Required** - Token for admin API access (set a strong secret!)

### **Database Configuration**

**Option 1: PostgreSQL URL (recommended)**

- `PG_URL` - Complete PostgreSQL URL (ex: `postgres://user:pass@localhost:5432/terrors?sslmode=disable`)

**Option 2: Individual variables**

- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - PostgreSQL user (default: postgres)
- `DB_PASS` - PostgreSQL password
- `DB_NAME` - Database name (default: terrors)

## 🎯 Features

_"Nobody trusts anybody now..." but you can trust this service_

### ✅ Core Features

- **Multi-Application Support** : Manage multiple applications from a single instance
- **Origin Validation** : Restrict error capture to specific domains
- **Error Persistence** : All errors stored in PostgreSQL with full details
- **Error Status Management** : Track errors as `new`, `treated`, or `deleted`
- **Fingerprinting** : Automatic error deduplication via SHA-1 fingerprinting
- **Statistics** : Error counts by status, type, and application

### ✅ Security Features

- **Admin Token Authentication** : All management endpoints protected
- **Origin Whitelist** : Prevent unauthorized apps from using your app_id
- **Token-based App Management** : Secure application creation and management

### ✅ Integration Features

- **Webhooks** : Discord and GitHub integrations (ready for implementation)
- **CQRS Dashboard** : Discover all available commands via `/overlook`
- **RESTful API** : Clean, predictable API design

### ✅ Developer Experience

- **Ultra-lightweight** : Single binary ~5-10MB
- **High performance** : Native compiled, instant startup
- **Hot reload** : Development with Air
- **Database migrations** : Versioned schema management
- **Horror-themed** : Every response comes with a scary movie quote

## 📖 Usage Guide

### Creating an Application

First, create an application to get an `app_id`:

```bash
curl -X POST http://localhost:3000/api/apps \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Awesome App",
    "description": "My application description",
    "origins": ["https://example.com", "https://app.example.com"]
  }'
```

Response:

```json
{
  "success": true,
  "message": "Application créée avec succès",
  "app": {
    "id": 1,
    "name": "My Awesome App",
    "appId": "app_a3f9k2m1",
    "description": "My application description",
    "origins": "https://example.com,https://app.example.com",
    "isActive": true,
    "createdAt": "2025-01-XX...",
    "updatedAt": "2025-01-XX..."
  },
  "dashboardToken": "a1b2c3d4e5f6...",
  "warning": "⚠️ Conservez ce token en sécurité, il ne sera affiché qu'une seule fois !"
```

**Note:** The `dashboardToken` is shown only once. Store it securely if you need it later.

### Using the Client Script

Add the script to your HTML:

```html
<script
  src="http://localhost:3000/cdn/terrors.js"
  data-app-id="app_a3f9k2m1"
></script>
```

The script will automatically:

- Capture JavaScript errors
- Capture unhandled promise rejections
- Compute error fingerprints
- Send errors to `/sadako` endpoint

### Managing Errors

List errors for an application:

```bash
curl -X GET "http://localhost:3000/api/errors?appId=app_a3f9k2m1&status=new&limit=50" \
  -H "Authorization: Bearer your-admin-token"
```

Update error status:

```bash
curl -X PATCH http://localhost:3000/api/errors/123 \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{"status": "treated"}'
```

Get error statistics:

```bash
curl -X GET "http://localhost:3000/api/errors/stats?appId=app_a3f9k2m1" \
  -H "Authorization: Bearer your-admin-token"
```

### Creating Webhooks

Webhooks are automatically triggered when a new error is captured. They run in the background and don't block error processing.

**Discord Webhook:**

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_a3f9k2m1" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "discord",
    "url": "https://discord.com/api/webhooks/123456789/abcdefghijklmnop",
    "config": {
      "username": "Terrors Bot"
    }
  }'
```

**GitHub Webhook** (creates issues automatically):

**Sans URL (recommandé)** - L'URL est construite automatiquement :

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_a3f9k2m1" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "github",
    "config": {
      "token": "ghp_xxxxxxxxxxxxxxxxxxxx",
      "owner": "username",
      "repo": "repository-name",
      "labels": ["bug", "error", "terrors"]
    }
  }'
```

**Avec URL complète** (optionnel) :

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_a3f9k2m1" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "github",
    "url": "https://api.github.com/repos/owner/repo/issues",
    "config": {
      "token": "ghp_xxxxxxxxxxxxxxxxxxxx",
      "owner": "username",
      "repo": "repository-name",
      "labels": ["bug", "error", "terrors"]
    }
  }'
```

**How it works:**

1. When an error is captured, it's saved to the database
2. All active webhooks for that app are automatically triggered
3. Each webhook sends a formatted notification:
   - **Discord**: Rich embed with error details
   - **GitHub**: Creates an issue with full error information

See `docs/webhooks.md` for detailed documentation.

### CQRS Dashboard

Discover all available commands:

```bash
curl http://localhost:3000/overlook
```

This returns a list of all available API commands with examples.

## 🎬 API Response Examples

_"Sweet dreams..."_

### Home endpoint response:

```json
{
  "message": "Welcome to the Overlook Hotel - Error monitoring service",
  "quote": "All work and no play makes Jack a dull boy",
  "year": "1980"
}
```

### Error capture response:

```json
{
  "success": true,
  "message": "Error captured and stored in the basement",
  "timestamp": "2025-08-22T01:05:00Z",
  "quote": "Here's Johnny!"
}
```

### Dashboard response:

```json
{
  "message": "Welcome to the Overlook Hotel - Command Center",
  "quote": "They're here",
  "commands": [
    {
      "name": "ListApps",
      "description": "Lister toutes les applications",
      "method": "GET",
      "path": "/api/apps",
      "example": {
        "headers": {
          "Authorization": "Bearer <ADMIN_TOKEN>"
        }
      }
    }
    // ... more commands
  ]
}
```

_"The quotes are randomly selected from classic horror movies (1970-2000)."_

## 🎭 Horror Movie References

This service pays homage to classic horror films with every response:

- **The Shining** (1980) - "Here's Johnny!", "All work and no play makes Jack a dull boy"
- **Halloween** (1978) - "The boogeyman is coming"
- **A Nightmare on Elm Street** (1984) - "Sweet dreams", "Don't fall asleep"
- **The Exorcist** (1973) - "The power of Christ compels you"
- **Alien** (1979) - "In space no one can hear you scream"
- **Scream** (1996) - "What's your favorite scary movie?"
- **Poltergeist** (1982) - "They're here"
- **The Ring** (2002) - "Seven days"
- And many more...

_"Seven days..." until your next error appears._

## 🚀 Development

_"Things start getting weird..." when you modify the code_

The project uses **Air** for hot reload. Modify a `.go` file and the server restarts automatically!

_"The server never sleeps, it just watches and waits..."_

### Database Migrations

Run migrations:

```bash
go run cmd/migrate/main.go
```

Migrations are automatically tracked and only run once.

## 🚢 CapRover Deployment

_"The call is coming from inside the server..."_

### Deploy to CapRover

1. **Push your code to Git**

```bash
git add .
git commit -m "🎭 Add horror-themed error monitoring service"
git push origin main
```

2. **Deploy via CapRover Dashboard**

   - Connect your Git repository
   - Set environment variables:
     - `PG_URL`: Your PostgreSQL connection string
     - `ADMIN_TOKEN`: Your admin token (keep it secret!)
     - `PORT`: 3000 (default)

3. **Or deploy via CLI**

```bash
caprover deploy --appName terrors --imageName terrors
```

### CapRover Features

- 🚢 **Auto-deployment** : Deploy on every Git push
- 🔒 **SSL/TLS** : Automatic HTTPS certificates
- 🔄 **Zero-downtime** : Rolling updates
- 📊 **Monitoring** : Built-in health checks
- 🌐 **Custom domains** : Easy domain configuration

_"The server never sleeps, it just watches and waits for errors in the cloud..."_

## 🧪 Testing

_"The call is coming from inside the tests..."_

### Unit Tests

Run Go unit tests:

```bash
# All tests
go test ./...

# Specific package
go test ./internal/services

# With coverage
go test -cover ./internal/services
```

**Configuration**: Set `PG_URL` or `PG_URL_TEST` environment variable.

### Manual Testing

Use the provided script to test all endpoints:

```bash
chmod +x test_manual.sh
./test_manual.sh
```

Or test manually with curl (see examples in [Usage Guide](#-usage-guide)).

**See `docs/testing.md` for a complete testing guide.**

## 🧪 Testing

_"The call is coming from inside the tests..."_

### Quick Start

**For beginners**: See `QUICK_START_TESTING.md` for a step-by-step guide.

### Unit Tests

```bash
# All tests
go test ./...

# Specific package with details
go test -v ./internal/services

# With coverage
go test -cover ./internal/services
```

**Configuration**: Set `PG_URL` or `PG_URL_TEST` environment variable.

### Manual Testing

```bash
# Linux/Mac
bash test_manual.sh

# Or with Makefile
make test-manual
```

**Documentation**:

- `QUICK_START_TESTING.md` - Guide rapide pour débutants
- `docs/testing.md` - Guide complet des tests

## 🔒 Security Notes

- **Admin Token** : Keep your `ADMIN_TOKEN` secret and use strong values
- **Origins** : Always configure allowed origins for production apps
- **HTTPS** : Use HTTPS in production to protect tokens and data
- **Database** : Use strong PostgreSQL passwords and restrict access

## 📝 License

_"The horror... the horror..."_

This project is open source. Use it, modify it, make it your own.

---

_"Welcome to the Overlook Hotel. We've been waiting for you..."_
