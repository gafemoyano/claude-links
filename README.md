# BE-Links - Dynamic Link Service

A Go + Fiber service that replaces Firebase Dynamic Links, providing smart redirects based on device platform and deep linking capabilities.

## 🚀 Quick Start

### Compile and Run

```bash
# 1. Start PostgreSQL database
docker compose up -d

# 2. Install dependencies and run
go mod download
go run main.go

# Or compile and run binary
go build -o be-links
./be-links
```

The server starts on `http://localhost:3000`

## 🛠 Technologies

- **Backend**: Go + Fiber
- **Database**: PostgreSQL
- **Containerization**: Docker Compose
- **Mobile Deep Links**: Universal Links (iOS) + App Links (Android)

## 📖 Detailed Setup

### Prerequisites

- Go 1.24+ installed
- Docker and Docker Compose installed

### 1. Clone and Setup

```bash
git clone <your-repo>
cd be-links
```

### 2. Start PostgreSQL Database

```bash
docker compose up -d
```

This starts PostgreSQL on port `5433` with:
- Database: `be-links`
- User: `felipe`
- Password: `password`

### 3. Run the Application

```bash
# Option 1: Using environment file
source .env && go run main.go

# Option 2: Direct environment variable
DATABASE_URL=postgres://felipe:password@localhost:5433/be-links?sslmode=disable go run main.go
```

The server will start on `http://localhost:3000`

### 4. Verify Setup

Test the health endpoint:
```bash
curl http://localhost:3000/health
```

Expected response:
```json
{"status":"ok"}
```

## 📖 API Usage

### Create a Short Link

```bash
curl -X POST http://localhost:3000/admin/create \
  -u admin:password \
  -H "Content-Type: application/json" \
  -d '{
    "deep_link": "myapp://product?id=987",
    "ios_store": "https://apps.apple.com/app/id123456789",
    "android_store": "https://play.google.com/store/apps/details?id=com.example.app",
    "title": "Product Page",
    "description": "Check out this awesome product"
  }'
```

Response:
```json
{
  "id": "abc123",
  "deep_link": "myapp://product?id=987",
  "ios_store": "https://apps.apple.com/app/id123456789",
  "android_store": "https://play.google.com/store/apps/details?id=com.example.app",
  "title": "Product Page",
  "description": "Check out this awesome product",
  "created_at": "2025-06-05T21:35:00Z",
  "updated_at": "2025-06-05T21:35:00Z",
  "click_count": 0
}
```

### Use the Short Link

Access the generated short link:
```bash
# This will redirect based on User-Agent
curl -L http://localhost:3000/abc123

# Test iOS redirect
curl -L -A "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X)" http://localhost:3000/abc123

# Test Android redirect  
curl -L -A "Mozilla/5.0 (Linux; Android 10; SM-G975F)" http://localhost:3000/abc123
```

## 🔗 How It Works

### Platform Detection

The service detects the user's platform via User-Agent and redirects accordingly:

- **iOS devices**: Redirects to App Store URL
- **Android devices**: Creates intent URI with fallback to Play Store
- **Unknown devices**: Returns JSON with all available links

### Redirect Logic

1. **iOS**: Direct redirect to `ios_store` URL
2. **Android**: Builds intent URI: `intent://path#Intent;scheme=myapp;package=com.example.app;S.browser_fallback_url=PLAY_STORE_URL;end`
3. **Other**: Returns JSON response with all links

## 📱 Mobile App Integration

### iOS Universal Links

1. Update `.well-known/apple-app-site-association` with your Team ID and app ID:
```json
{
  "applinks": {
    "apps": [],
    "details": [
      {
        "appID": "YOUR_TEAM_ID.com.yourapp.bundle",
        "paths": ["*"]
      }
    ]
  }
}
```

2. Add Associated Domains capability in Xcode:
   - `applinks:yourdomain.com`

### Android App Links

1. Update `.well-known/assetlinks.json` with your package and SHA256 fingerprint:
```json
[
  {
    "relation": ["delegate_permission/common.handle_all_urls"],
    "target": {
      "namespace": "android_app",
      "package_name": "com.yourapp.package",
      "sha256_cert_fingerprints": [
        "YOUR_SHA256_FINGERPRINT"
      ]
    }
  }
]
```

2. Add intent filters to `android/app/src/main/AndroidManifest.xml`:
```xml
<intent-filter android:autoVerify="true">
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data android:scheme="https" android:host="yourdomain.com" />
</intent-filter>
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file or set these environment variables:

```bash
# Database
DATABASE_URL=postgres://felipe:password@localhost:5433/be-links?sslmode=disable

# Server
PORT=3000

# Admin Auth (change in production!)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=password
```

### Database Schema

The application automatically creates this table:

```sql
CREATE TABLE links (
    id VARCHAR(255) PRIMARY KEY,
    deep_link TEXT NOT NULL,
    ios_store TEXT NOT NULL,
    android_store TEXT NOT NULL,
    title TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    click_count INTEGER DEFAULT 0
);
```

## 🏗 Project Structure

```
be-links/
├── main.go                              # Application entry point
├── go.mod                               # Go dependencies
├── docker-compose.yml                   # PostgreSQL setup
├── .env                                 # Environment variables
├── routes/
│   ├── redirect.go                      # Short link redirect handling
│   └── admin.go                         # Protected admin routes
├── handlers/
│   └── link_handler.go                  # Business logic
├── models/
│   └── link.go                          # Data models
├── storage/
│   └── db.go                            # Database operations
└── .well-known/
    ├── apple-app-site-association       # iOS Universal Links
    └── assetlinks.json                  # Android App Links
```

## 🚦 API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/health` | Health check | None |
| GET | `/:shortcode` | Redirect to app/store | None |
| POST | `/admin/create` | Create short link | Basic Auth |
| GET | `/.well-known/apple-app-site-association` | iOS Universal Links | None |
| GET | `/.well-known/assetlinks.json` | Android App Links | None |

## 🔒 Security

- Admin endpoints protected with Basic Authentication
- Input validation on link creation
- SQL injection prevention with parameterized queries
- No open redirects (only to predefined store URLs)

## 📊 Analytics

The service tracks click counts for each short link. Access the database directly to view analytics:

```bash
psql postgres://felipe:password@localhost:5433/be-links -c "SELECT id, click_count, created_at FROM links ORDER BY click_count DESC;"
```

## 🚀 Production Deployment

### Google Cloud Run

1. Build and push Docker image:
```bash
docker build -t gcr.io/PROJECT_ID/be-links .
docker push gcr.io/PROJECT_ID/be-links
```

2. Deploy to Cloud Run:
```bash
gcloud run deploy be-links \
  --image gcr.io/PROJECT_ID/be-links \
  --platform managed \
  --region us-central1 \
  --set-env-vars DATABASE_URL=your-production-db-url
```

### Environment Variables for Production

- Set strong `ADMIN_USERNAME` and `ADMIN_PASSWORD`
- Use a production PostgreSQL instance
- Enable HTTPS (required for Universal Links)

## 🔧 Development

### Running Tests

```bash
go test ./...
```

### Database Management

```bash
# Access database
psql postgres://felipe:password@localhost:5433/be-links

# View all links
SELECT * FROM links;

# Reset click counts
UPDATE links SET click_count = 0;
```

### Stopping Services

```bash
# Stop PostgreSQL container
docker compose down

# Remove volumes (deletes data)
docker compose down -v
```

## 📝 License

MIT License - see LICENSE file for details.