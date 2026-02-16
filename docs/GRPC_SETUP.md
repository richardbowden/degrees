# gRPC + gRPC-Gateway Setup Complete

## What Was Created

### 1. Protocol Buffer Definitions

**`proto/p402/v1/common.proto`**
- User message type
- Error message type
- Pagination types

**`proto/p402/v1/user_service.proto`**
- UserService with 6 endpoints:
  - GetUser - Get user by ID
  - UpdateUser - Update user profile
  - EnableUser - Enable user account (admin)
  - DisableUser - Disable user account (admin)
  - SetUserSysop - Set sysop status (sysop)
  - ListUsers - List all users (admin)

**`proto/p402/v1/auth_service.proto`**
- AuthService with 6 endpoints:
  - Register - User registration
  - VerifyEmail - Email verification
  - Login - User login
  - Logout - User logout
  - ChangePassword - Change password
  - ResetPassword - Request password reset

### 2. Configuration Files

**`proto/buf.yaml`**
- Buf configuration with linting rules
- Google APIs dependency for HTTP annotations

**`proto/buf.gen.yaml`**
- Code generation for:
  - Go protobuf messages
  - gRPC server/client stubs
  - gRPC-Gateway reverse proxy
  - OpenAPI/Swagger documentation

### 3. Generated Code

**`internal/pb/p402/v1/*.pb.go`** (7 files)
- Protobuf message types
- gRPC service definitions
- HTTP gateway handlers
- All auto-generated, never edit manually

**`docs/openapi/api.swagger.json`**
- Complete OpenAPI spec for REST API
- Can be used with Swagger UI, Postman, etc.

### 4. Service Implementations

**`internal/grpc/user_service.go`**
- UserServiceServer implementation
- Connects to existing `services.UserService`
- EnableUser, DisableUser, ListUsers implemented
- GetUser, UpdateUser, SetUserSysop stubbed (TODO)

**`internal/grpc/auth_service.go`**
- AuthServiceServer implementation stub
- All methods return Unimplemented
- Ready for implementation

### 5. Build Integration

**`magefiles/build.go`**
- Added `protoGen()` function
- Integrated into `Gen()` target
- Runs before SQL generation

## Architecture

```
┌─────────────┐
│   Clients   │
│ iOS/Web/Go  │
└──────┬──────┘
       │
       ├─────────────┬──────────────┐
       │             │              │
   gRPC (9090)   REST (8080)    Native Go
       │             │              │
       │      ┌──────▼──────┐      │
       │      │  Gateway    │      │
       │      │  (gRPC-Web) │      │
       │      └──────┬──────┘      │
       │             │              │
       └─────────────┴──────────────┘
                     │
            ┌────────▼────────┐
            │  gRPC Services  │
            │ UserService     │
            │ AuthService     │
            └────────┬────────┘
                     │
            ┌────────▼────────┐
            │ Business Logic  │
            │ services/*      │
            └────────┬────────┘
                     │
            ┌────────▼────────┐
            │  Repositories   │
            │ repos/*         │
            └────────┬────────┘
                     │
            ┌────────▼────────┐
            │    Database     │
            └─────────────────┘
```

## Next Steps

### Phase 1: Wire Up Server (HIGH PRIORITY)

**File**: `cmd/server/server_run_cmd.go`

Add gRPC server startup:

```go
import (
    "net"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

    pb "github.com/typewriterco/p402/internal/pb/p402/v1"
    grpcsvr "github.com/typewriterco/p402/internal/grpc"
)

func serverRun(c *cli.Context) error {
    // ... existing setup ...

    // Create gRPC server
    grpcServer := grpc.NewServer()

    // Register services
    userGrpcSvc := grpcsvr.NewUserServiceServer(userSvc)
    pb.RegisterUserServiceServer(grpcServer, userGrpcSvc)

    authGrpcSvc := grpcsvr.NewAuthServiceServer(authNService, signUpSvc)
    pb.RegisterAuthServiceServer(grpcServer, authGrpcSvc)

    // Enable reflection (for grpcurl/grpcui)
    reflection.Register(grpcServer)

    // Start gRPC server
    grpcListener, _ := net.Listen("tcp", ":9090")
    go func() {
        log.Info().Msg("gRPC server listening on :9090")
        grpcServer.Serve(grpcListener)
    }()

    // Create gRPC-Gateway
    ctx := context.Background()
    gwmux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

    pb.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, "localhost:9090", opts)
    pb.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, "localhost:9090", opts)

    // Mount gateway
    r := chi.NewRouter()
    r.Mount("/", gwmux)

    // Start HTTP gateway
    httpServer := &http.Server{Addr: ":8080", Handler: r}
    log.Info().Msg("HTTP gateway listening on :8080")
    return httpServer.ListenAndServe()
}
```

### Phase 2: Add Authentication Middleware

Create gRPC interceptor for auth:

```go
// internal/grpc/auth_interceptor.go
func AuthInterceptor(authSvc *services.AuthNService) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // Extract session token from metadata
        // Validate session
        // Add user to context
        return handler(ctx, req)
    }
}
```

Register interceptor:
```go
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(AuthInterceptor(authNService)),
)
```

### Phase 3: Implement Remaining Endpoints

**Priority order:**
1. **AuthService.Login** - Most critical
2. **AuthService.Register** - User signup
3. **AuthService.VerifyEmail** - Complete signup flow
4. **UserService.GetUser** - Profile viewing
5. **AuthService.ChangePassword** - Security
6. **UserService.UpdateUser** - Profile editing

### Phase 4: Add Validation

Use buf validate for proto-level validation:

```protobuf
import "buf/validate/validate.proto";

message RegisterRequest {
  string email = 1 [(buf.validate.field).string.email = true];
  string password = 2 [(buf.validate.field).string.min_len = 8];
}
```

### Phase 5: Client Implementation

**iOS App**:
```bash
# Generate Swift code
cd proto
buf generate --template buf.gen.swift.yaml
```

**Web Frontend**:
```javascript
// Use REST endpoints via gateway
const response = await fetch('/api/v1/admin/users');
const data = await response.json();
```

## Testing

### Test gRPC directly

```bash
# Install grpcurl
brew install grpcurl

# List services
grpcurl -plaintext localhost:9090 list

# Call endpoint
grpcurl -plaintext \
  -d '{"user_id": 1}' \
  localhost:9090 \
  p402.v1.UserService/EnableUser
```

### Test REST endpoints

```bash
curl http://localhost:8080/api/v1/admin/users
```

### Interactive UI

```bash
# Install grpcui
brew install grpcui

# Open interactive UI
grpcui -plaintext localhost:9090
```

## Benefits Achieved

✅ **Single source of truth** - Proto files define everything
✅ **Auto-generated OpenAPI** - `docs/openapi/api.swagger.json`
✅ **Type-safe clients** - For Go, Swift, JS, etc.
✅ **Both gRPC and REST** - Same code, both protocols
✅ **Contract enforcement** - Breaking changes detected by buf
✅ **Integrated build** - `mage` generates everything

## Files Changed

- `proto/` (new) - All proto definitions
- `internal/pb/` (new) - Generated code
- `internal/grpc/` (new) - Service implementations
- `docs/openapi/` (new) - OpenAPI spec
- `magefiles/build.go` - Added proto generation

## Documentation

- `proto/README.md` - Proto development guide
- `docs/openapi/api.swagger.json` - API documentation
- This file - Setup and next steps
