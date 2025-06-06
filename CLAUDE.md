# CLAUDE.md

## Project: Dynamic Link Service with Go + Fiber

### 🧩 Overview

This service replaces Firebase Dynamic Links. It supports generating short URLs that redirect users to appropriate destinations based on device platform, deep linking into a Flutter app, or fallback to the App/Play Store when the app isn't installed.

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
- [x] Redirect:
  - iOS → universal link or App Store
  - Android → intent URI or Play Store
- [x] Handle fallback logic
- [x] Host Apple + Android association files
- [ ] (Optional) Analytics/logging
- [ ] (Optional) Admin auth for link creation

---

## 🔗 Short Link Format

Example:  
`https://yourdomain.com/abc123`

Maps to:
```json
{
  "id": "abc123",
  "deep_link": "myapp://product?id=987",
  "ios_store": "https://apps.apple.com/app/id123456789",
  "android_store": "https://play.google.com/store/apps/details?id=com.example.app"
}
```

---

## 🚦 Route Behavior

### `GET /:shortcode`

- Parse `User-Agent` to detect device
- If app installed:
  - iOS: redirect via universal link (e.g., `https://yourdomain.com/app/path`)
  - Android: use `intent://` scheme
- If app not installed:
  - Redirect to respective app store
- Fallback: Show a landing page or error

### `POST /create` (Protected)

- JSON body with deep link, app store URLs
- Generates random short code
- Stores in PostgreSQL

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

## 📱 Flutter Deep Link Handling

### iOS

- Configure `applinks:yourdomain.com` in Xcode
- Add Associated Domains entitlement
- Host `.well-known/apple-app-site-association`

### Android

- Add intent filters to `AndroidManifest.xml`
- Host `.well-known/assetlinks.json`
- Sign APK with matching SHA256 cert

### Flutter Setup

Use `uni_links` or `go_router`:
```dart
void main() async {
  final initialUri = await getInitialUri();
  runApp(MyApp(initialUri: initialUri));
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
- Link expiration
- QR code generator
- Admin UI with auth

---

Made for modern app teams migrating from Firebase Dynamic Links.
