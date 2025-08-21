# 🎭 Terrors - Error Monitoring Service

> _"The call is coming from inside the house..."_

A horror-themed JavaScript error monitoring service written in Go. When your application crashes, we'll be there to capture the screams... I mean, the errors.

_"In space no one can hear you scream, but in our logs, every error is preserved forever."_

## 🏚️ Project Structure

_"They're here..." - The files are watching you_

```
terrors/
├── cmd/
│   └── server/
│       └── main.go          # The main entrance to our haunted mansion
├── internal/
│   ├── api/
│   │   └── handlers.go      # The handlers that never sleep
│   ├── database/
│   │   └── db.go           # The basement where errors are stored
│   └── models/
│       └── error.go        # The ghosts of your past mistakes
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

3. **Install Air (hot reload)**

```bash
go install github.com/air-verse/air@latest
```

4. **Configure environment variables**

Copy the example file:

```bash
cp env.example .env
```

Then modify `.env` according to your configuration.

5. **Start the server**

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

- `GET /` - Welcome to the Overlook Hotel
- `POST /sadako` - The call is coming from inside the house (error endpoint)
- `GET /cdn/terrors.js` - The script that haunts your browser

## 👻 Client Script Usage

_"Don't fall asleep..."_

```html
<script src="http://localhost:3000/cdn/terrors.js" appid="your-app-id"></script>
```

_"The script will watch your application and capture every error that dares to appear."_

## 🔮 Environment Variables

_"The power of Christ compels you..." to configure these properly_

Copy `env.example` to `.env` and configure:

### **Option 1: PostgreSQL URL (recommended)**

- `PG_URL` - Complete PostgreSQL URL (ex: `postgres://user:pass@localhost:5432/terrors`)

### **Option 2: Individual variables**

- `PORT` - Server port (default: 3000)
- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - PostgreSQL user (default: postgres)
- `DB_PASS` - PostgreSQL password
- `DB_NAME` - Database name (default: terrors)

## 🚀 Development

_"Things start getting weird..." when you modify the code_

The project uses **Air** for hot reload. Modify a `.go` file and the server restarts automatically!

_"The server never sleeps, it just watches and waits..."_

## 🎯 Advantages of the Go Version

_"Nobody trusts anybody now..." but you can trust this service_

- ✅ **Ultra-lightweight** : Single binary ~5-10MB
- ✅ **High performance** : Native compiled, instant startup
- ✅ **Simple** : No external frameworks
- ✅ **Standard** : Conventional Go structure
- ✅ **Easy deployment** : Portable binary
- ✅ **Horror-themed** : Every error comes with a scary movie quote

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

_"The quotes are randomly selected from classic horror movies (1970-2000)."_

## 🎭 Horror Movie References

This service pays homage to classic horror films with every response:

- **The Shining** (1980) - "Here's Johnny!"
- **Halloween** (1978) - "The boogeyman is coming"
- **A Nightmare on Elm Street** (1984) - "Sweet dreams"
- **The Exorcist** (1973) - "The power of Christ compels you"
- **Alien** (1979) - "In space no one can hear you scream"
- **Scream** (1996) - "What's your favorite scary movie?"
- And many more...

_"Seven days..." until your next error appears._

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
