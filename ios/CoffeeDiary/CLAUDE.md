# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

iOS client for Coffee Diary — an espresso brewing tracker. Connects to a Go backend at `coffee.oglimmer.com` with OIDC authentication (Keycloak). Users log brewing parameters, manage beans/equipment, and track tasting notes.

## Build & Run

```bash
# Build from command line
xcodebuild -project CoffeeDiary.xcodeproj -scheme CoffeeDiary -sdk iphonesimulator -configuration Debug build

# Run tests (when added)
xcodebuild -project CoffeeDiary.xcodeproj -scheme CoffeeDiary -sdk iphonesimulator -configuration Debug test
```

Open `CoffeeDiary.xcodeproj` in Xcode for development, previews, and running on simulator/device.

## Project Configuration

- **Xcode 16+** with `PBXFileSystemSynchronizedRootGroup` — files added to `CoffeeDiary/` are auto-discovered, no pbxproj editing needed
- **Bundle ID:** `com.oglimmer.CoffeeDiary`
- **Deployment target:** iOS 26.4 (set in project, targeting iOS 17+ at runtime)
- **Swift 5** with `SWIFT_DEFAULT_ACTOR_ISOLATION = MainActor` and `SWIFT_APPROACHABLE_CONCURRENCY = YES`
- **No SPM packages or CocoaPods** currently — pure Apple frameworks only
- Portrait orientation on iPhone, all orientations on iPad

## Architecture

- **Pattern:** MVVM with SwiftUI
- **UI Framework:** SwiftUI (Composition API style with `@Observable` / `@State`)
- **Networking:** URLSession with cookie-based session auth (session cookie from OIDC flow)
- **Auth:** OIDC via `ASWebAuthenticationSession` → Keycloak login → session cookie

## Backend API

The Go backend lives at `/Users/oli/dev/coffee-diary/backend`. Production URL: `https://coffee.oglimmer.com`.

Key endpoints (all under `/api`):
- `GET /api/auth/login` — redirects to Keycloak OIDC login
- `GET /api/auth/callback` — OIDC callback, sets session cookie
- `GET /api/auth/me` — returns `{id, username}` or 401
- `GET /api/auth/logout` — clears session
- `GET/POST/PUT/DELETE /api/diary-entries[/{id}]` — paginated (`page`, `size`), filterable (`coffeeId`, `sieveId`, `dateFrom`, `dateTo`, `ratingMin`), sortable (`sort=field,dir`)
- `GET/POST/DELETE /api/coffees[/{id}]` — coffee beans
- `GET/POST/DELETE /api/sieves[/{id}]` — sieves/filters

DateTime format: `2006-01-02T15:04:05` (no timezone). Rating: 1-5. Default temperature: 93°C.
