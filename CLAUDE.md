# CLAUDE.md

## Project: Universal/App Link Service with Go + Fiber

### 🧩 Overview

This service replaces Firebase Dynamic Links using modern universal links (iOS) and app links (Android). It generates short URLs that redirect users to appropriate destinations based on device platform, using HTTPS-based universal/app links for seamless deep linking into Flutter apps, with fallback to App/Play Store when the app isn't installed.

---

## 🛠 Technologies

- **Backend**: Go + Fiber
- **Database**: PostgreSQL
- **Mobile App**: Flutter (currently using Firebase Dynamic Links)
- **Hosting**: Google Cloud Run (HTTPS-enabled)

---

## 📦 Project Structure

```
.
├── main.go
├── routes/
│   └── redirect.go
│   └── admin.go
├── handlers/
│   └── link_handler.go
├── models/
│   └── link.go
├── storage/
│   └── db.go
├── .well-known/
│   └── apple-app-site-association
│   └── assetlinks.json
├── CLAUDE.md
└── go.mod
```

---

## 🎯 Features

- [x] Create short links for app content  
- [x] Detect platform via User-Agent
- [x] Redirect using modern approach:
  - iOS → Universal Links (HTTPS URLs) with App Store fallback
  - Android → App Links (HTTPS URLs) with Play Store fallback
- [x] Handle fallback logic when app not installed
- [x] Host Apple App Site Association and Android Asset Links files
- [x] Analytics/logging with click tracking
- [x] Admin auth for link creation

---

## 🔗 Short Link Format

Example:
`https://yourdomain.com/abc123`

Maps to:

```json
{
  "id": "abc123",
  "universal_link": "https://yourdomain.com/app/product?id=987",
  "ios_store": "https://apps.apple.com/co/app/trii/id1513826307",
  "android_store": "https://play.google.com/store/apps/details?id=com.triico.app&hl=en",
  "title": "Product Page",
  "description": "Check out this awesome product"
}
```

---

## 🚦 Route Behavior

### `GET /:shortcode`

- Parse `User-Agent` to detect device platform
- Smart redirect behavior:
  - **iOS**: Redirect to universal link (HTTPS URL that opens app if installed, otherwise goes to App Store)
  - **Android**: Redirect to app link (HTTPS URL that opens app if installed, otherwise goes to Play Store)  
  - **Desktop/Unknown**: Show landing page with all options
- All redirects use HTTPS URLs for maximum compatibility

### `POST /admin/create` (Protected)

- JSON body with universal link path, app store URLs, title, description
- Generates random short code
- Stores in PostgreSQL with click tracking

---

## 🔐 Security Considerations

- Secure link creation with basic auth or API key
- Validate all incoming data (especially redirect targets)
- Prevent open redirects

---

## 🧪 Testing Plan

- Cold start / warm start on iOS & Android
- Deep link handling via `uni_links` or `go_router`
- Play Store and App Store fallback
- Universal Links and AssetLinks validation
- Redirect latency & correctness

---

## 📱 Flutter Universal/App Link Handling

### iOS Universal Links

- Configure `applinks:yourdomain.com` in Associated Domains capability
- Add entitlement in Xcode project settings
- Host `.well-known/apple-app-site-association` at domain root
- Links use HTTPS format: `https://yourdomain.com/app/path`

### Android App Links

- Add intent filters with `android:autoVerify="true"` to `AndroidManifest.xml`
- Host `.well-known/assetlinks.json` at domain root
- Sign APK with matching SHA256 certificate fingerprint
- Links use HTTPS format: `https://yourdomain.com/app/path`

### Flutter Setup

Use `go_router` for modern routing:

```dart
final GoRouter _router = GoRouter(
  routes: [
    GoRoute(
      path: '/app/:action',
      builder: (context, state) {
        final action = state.pathParameters['action'];
        return handleDeepLink(action, state.uri.queryParameters);
      },
    ),
  ],
);

void main() {
  runApp(MaterialApp.router(routerConfig: _router));
}
```

---

## 🌐 Hosting Requirements

- HTTPS mandatory (iOS requirement for universal links)
- Serve `.well-known` files at root with `Content-Type: application/json`
- No redirects for AASA or assetlinks.json

---

## 🚀 Deployment Checklist

- [ ] Serve `.well-known` files over HTTPS
- [ ] Set up PostgreSQL and migrations
- [ ] Configure mobile apps for universal/app links
- [ ] Protect admin endpoints
- [ ] Test all platforms

---

## 🧠 Future Improvements

- Analytics dashboard (clicks, devices, countries)
- Link expiration and custom domains
- QR code generator for links
- Admin UI with enhanced auth
- A/B testing for different fallback strategies

---

Made for modern app teams migrating from Firebase Dynamic Links to universal/app links.
