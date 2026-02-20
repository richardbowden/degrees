# 40 Degrees Car Detailing API

## Project Overview
Backend API for 40 Degrees Car Detailing, a premium mobile detailing business
based in Perth, Western Australia. Uses Bowden's Own premium Australian products.
Environmental focus: wildlife and marine life safe products and practices.

Built on the degrees base service template. Being renamed to github.com/typewriterco/40degrees.

## Tech Stack
- Go with gRPC services and HTTP/REST gateway proxy
- PostgreSQL with sqlc (app schema + fga + river)
- OpenFGA for authorization (fga schema)
- River for background jobs (river schema)
- Template engine with versioning (from degrees base)
- Notification system (from degrees base)
- Stripe for payments (AUD, cents as BIGINT)
- golang-migrate for migrations
- urfave/cli for CLI structure

## Architecture
This follows the degrees architecture pattern:
- gRPC service definitions (protobuf) with HTTP/REST gateway proxy
- Clients consume a REST API, which the gateway translates to gRPC calls internally
- `internal/services/` — Business logic, injected with repos
- `internal/repos/` — Data access layer (sqlc generated)
- `internal/transport/` — gRPC server implementations and HTTP gateway
- `internal/jobs/` — River background job definitions and handlers
- `cmd/server/` — CLI entry point, config loading, dependency wiring

## CRITICAL: This is Built on degrees
- ALL existing code, services, handlers, and config MUST be preserved
- User registration, login, sessions, OpenFGA, River, templates, notifications
  are all working and must not be removed or rewritten
- New detailing features are ADDED alongside existing degrees code
- New database tables go in the PUBLIC schema alongside existing degrees tables
- New migrations are APPENDED to the existing migration sequence
- New services, repos, and gRPC servers are ADDED, existing ones untouched
- New proto definitions are ADDED, existing ones untouched
- Server startup wiring is EXTENDED to include new services, not rewritten

## Patterns
- gRPC services define the API contract via protobuf
- HTTP gateway proxies REST requests to gRPC
- Services receive repos via constructor injection
- OpenFGA checks happen in services, not transport layer
- All money stored as BIGINT cents, never floats
- Use RFC 7807 Problem Details for errors (problems package)
- Database sessions (already working from degrees)
- Use the existing template engine for all outbound emails and notifications
- River jobs should use the notification service for delivery

## Database Schemas
- The app schema (named after the project, renamed from degrees during setup)
  holds all application tables, including new detailing tables
- fga: OpenFGA authorization tables (DO NOT TOUCH)
- river: River job queue tables (DO NOT TOUCH)

## Domain: Car Detailing (NEW, added to degrees)
- Services have categories, base prices, duration estimates, and optional add-ons
- Customers have profiles linked to existing users table, with vehicles
- Bookings require a 30% deposit via Stripe, balance collected on service day
- Scheduling is duration-aware with configurable business hours and blackout dates
- Service history tracks completed work with staff notes, products used, and photos
- Notes can be internal (staff only) or visible to the customer

## Business Rules
- Deposits are 30% of service total (configurable per service)
- Booking requires minimum 24-hour advance notice
- Business hours: Mon-Sat 7:00 AM to 5:00 PM AWST
- 30-minute buffer between bookings
- Cancelled bookings within 24 hours may forfeit deposit
- Cart supports guest sessions (token-based) and authenticated users

## Build and Run
```
go build -o 40degrees ./cmd/server
./40degrees db migration up
./40degrees server run
```

## Environment Variables
```
DATABASE_URL=postgres://degrees:letmein@localhost:5432/40degrees
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
SMTP_HOST=localhost
SMTP_PORT=1025
BASE_URL=http://localhost:3000
```
