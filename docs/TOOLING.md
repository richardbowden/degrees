# Development Tooling

## Zero-Installation Tools ✅

These tools work via `go run` - no installation required:

### buf (Protocol Buffer Tooling)
Automatically invoked via mage. Used for:
- Generating Go code from .proto files
- Linting proto files
- Breaking change detection
- Formatting

### How It Works
```bash
# This happens automatically in mage:
go run github.com/bufbuild/buf/cmd/buf@latest generate

# Downloads buf on first run
# Caches for subsequent runs
# Specified version in go.mod via tools.go
```

## Mage Commands

### Build & Generate
```bash
mage              # Default: runs proto gen, sql gen, go gen, builds debug binary
mage debug        # Same as default
mage release      # Production build with tests
mage gen          # Run all code generation (proto + sql + go)
mage clean        # Clean build artifacts
```

### Proto Commands
```bash
mage proto:lint      # Lint .proto files
mage proto:format    # Auto-format .proto files
mage proto:breaking  # Check for breaking API changes vs main branch
```

### Testing Commands
```bash
mage test         # Run Go tests
mage run          # Build and run server with human-readable logs
```

### gRPC Testing (requires running server)
```bash
mage grpc:list    # List all available gRPC services
mage grpc:ui      # Open interactive gRPC testing UI at http://localhost:8081
```

## Optional Tools (Install Once)

### grpcurl - Command-line gRPC client
```bash
# Install
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Usage
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext localhost:9090 p402.v1.UserService/ListUsers
grpcurl -plaintext -d '{"user_id": 1}' localhost:9090 p402.v1.UserService/EnableUser

# With auth
grpcurl -plaintext \
  -H "Authorization: Bearer <token>" \
  -d '{"user_id": 1}' \
  localhost:9090 p402.v1.UserService/EnableUser
```

### grpcui - Interactive gRPC web UI
```bash
# Install
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest

# Usage (or use mage grpc:ui)
grpcui -plaintext localhost:9090
# Opens browser at http://localhost:8080
```

### evans - Alternative gRPC client
```bash
# Install
go install github.com/ktr0731/evans@latest

# Usage (interactive REPL)
evans -r repl -p 9090
# Then type: service p402.v1.UserService
# Then type: call ListUsers
```

## Keeping Tools Updated

### Option 1: tools.go (Recommended)
We've created `tools/tools.go` that tracks tool versions in `go.mod`:

```bash
# Add tools to go.mod
go mod tidy

# Install a specific tool
go install github.com/fullstorydev/grpcurl/cmd/grpcurl

# Run without installing
go run github.com/fullstorydev/grpcurl/cmd/grpcurl@latest -plaintext localhost:9090 list
```

### Option 2: Update Tools
```bash
# Update all tools
go get -u github.com/bufbuild/buf/cmd/buf@latest
go get -u github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go get -u github.com/fullstorydev/grpcui/cmd/grpcui@latest

go mod tidy
```

## IDE Integration

### VS Code
Install extensions:
- **vscode-proto3** - Proto syntax highlighting
- **Buf** - Official buf extension
- **gRPC Protobuf Support** - Enhanced proto support

Settings (`.vscode/settings.json`):
```json
{
  "protoc": {
    "path": "buf",
    "compile_on_save": false
  },
  "buf.lintOnSave": true
}
```

### GoLand / IntelliJ
1. Enable Protocol Buffers plugin (built-in)
2. Set buf as the protoc compiler
3. Configure proto import paths to include `proto/`

## Workflow

### Daily Development
```bash
# 1. Make changes to .proto files
vim proto/p402/v1/user_service.proto

# 2. Generate code (happens automatically on build)
mage

# 3. Run server
mage run

# 4. Test in another terminal
mage grpc:ui
# OR
curl http://localhost:8080/api/v1/admin/users
```

### Before Committing
```bash
# 1. Lint protos
mage proto:lint

# 2. Format protos
mage proto:format

# 3. Check for breaking changes
mage proto:breaking

# 4. Run tests
mage test

# 5. Build release
mage release
```

### Adding New Proto Service
```bash
# 1. Create proto file
vim proto/p402/v1/project_service.proto

# 2. Define service with HTTP annotations
# 3. Generate code
mage

# 4. Implement in internal/grpc/
vim internal/grpc/project_service.go

# 5. Register in server_run_cmd.go
```

## Tool Dependencies Summary

| Tool | Required | Installation | Purpose |
|------|----------|--------------|---------|
| **buf** | ✅ Yes | Auto via `go run` | Proto compilation |
| **mage** | ✅ Yes | `go install github.com/magefile/mage@latest` | Build system |
| **sqlc** | ✅ Yes | Auto via `go tool sqlc` | SQL generation |
| **grpcurl** | ⚪ Optional | `go install` | CLI testing |
| **grpcui** | ⚪ Optional | `go install` | UI testing |
| **evans** | ⚪ Optional | `go install` | Alternative client |

## Troubleshooting

### buf not found
```bash
# buf is invoked via go run, should work automatically
# If issues, clear module cache:
go clean -modcache
```

### Generated code import errors
```bash
# Regenerate everything
mage clean
mage
```

### gRPC service not found
```bash
# Check server is running
lsof -i :9090

# Check reflection is enabled
grpcurl -plaintext localhost:9090 list

# Verify service registration in server_run_cmd.go
```

### Breaking changes detected
```bash
# View what changed
mage proto:breaking

# If intentional (new major version):
# 1. Create proto/p402/v2/ directory
# 2. Copy and modify protos there
# 3. Keep v1 for backward compatibility
```

## No External Dependencies Needed! 🎉

The key insight: We use `go run <package>@latest` for everything, which:
- Downloads tools on demand
- Caches them in Go's module cache
- Uses specific versions from go.mod
- No global installs required
- Works in CI/CD without setup

Only install grpcurl/grpcui if you want them for manual testing.
