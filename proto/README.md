# gRPC Protocol Definitions

This directory contains the Protocol Buffer definitions for the p402 API.

## Structure

```
proto/
├── buf.yaml              # Buf configuration
├── buf.gen.yaml          # Code generation config
└── p402/v1/             # API version 1
    ├── common.proto      # Shared types (User, Error, Pagination)
    ├── user_service.proto    # User management endpoints
    └── auth_service.proto    # Authentication endpoints
```

## Generated Code

Generated files are placed in:
- **Go Code**: `internal/pb/p402/v1/*.pb.go` - Protobuf messages
- **gRPC Code**: `internal/pb/p402/v1/*_grpc.pb.go` - gRPC server/client stubs
- **Gateway Code**: `internal/pb/p402/v1/*.pb.gw.go` - HTTP reverse proxy
- **OpenAPI**: `docs/openapi/api.swagger.json` - OpenAPI/Swagger spec

## Usage

### Generate Code

```bash
# Using mage (recommended)
mage

# Or manually
cd proto && go run github.com/bufbuild/buf/cmd/buf@latest generate
```

### Lint Protos

```bash
cd proto && go run github.com/bufbuild/buf/cmd/buf@latest lint
```

### Check for Breaking Changes

```bash
cd proto && go run github.com/bufbuild/buf/cmd/buf@latest breaking --against '.git#branch=main'
```

### Format Protos

```bash
cd proto && go run github.com/bufbuild/buf/cmd/buf@latest format -w
```

## Adding New Services

1. Create new `.proto` file in `p402/v1/`
2. Define service and messages
3. Add HTTP annotations for REST mapping:
   ```protobuf
   rpc GetUser(GetUserRequest) returns (GetUserResponse) {
     option (google.api.http) = {
       get: "/api/v1/user/{user_id}"
     };
   }
   ```
4. Run `mage` to generate code
5. Implement service in `internal/grpc/`

## Client Usage

### iOS/macOS (Swift)

```swift
import GRPC

let client = P402_V1_UserServiceClient(channel: channel)
let request = P402_V1_ListUsersRequest()
let response = try await client.listUsers(request)
```

### Web (REST via Gateway)

```javascript
fetch('/api/v1/admin/users')
  .then(response => response.json())
```

### Go (gRPC)

```go
import pb "github.com/typewriterco/p402/internal/pb/p402/v1"

client := pb.NewUserServiceClient(conn)
resp, _ := client.ListUsers(ctx, &pb.ListUsersRequest{})
```

## Dependencies

- **buf.build/googleapis/googleapis** - Google API annotations for HTTP mapping

## Best Practices

1. **Never break the API** - Use buf breaking checks
2. **Version everything** - New versions go in `p402/v2/`
3. **Document fields** - Add comments to proto messages
4. **Use pagination** - For list endpoints, use `PaginationRequest`/`PaginationResponse`
5. **Proper HTTP verbs** - GET for reads, POST for writes, PUT for updates, DELETE for deletes
