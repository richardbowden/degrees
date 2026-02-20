# Internal Code

## CRITICAL
- DO NOT modify, delete, or rewrite any existing degrees code
- ADD new files alongside existing ones
- Extend server wiring to register new services, do not replace it

## Architecture
- services/ contains business logic. Services are injected with repos.
- repos/ contains data access. Use sqlc generated code.
- transport/grpc/ contains gRPC server implementations injected with services.
- transport/gateway/ contains the HTTP/REST proxy (from degrees base).
- jobs/ contains River background job definitions and handlers.

## Rules
- gRPC servers never access the database directly, always through services
- Services never import transport layer packages
- OpenFGA authorization checks happen in the service layer
- All money is BIGINT cents, never use float for money
- Use the problems package for RFC 7807 error responses
- Use the existing notification service and template engine for emails
- River jobs use the notification service for delivery
- Follow existing degrees constructor injection patterns
