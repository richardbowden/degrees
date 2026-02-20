# Proto Definitions

## Rules
- DO NOT modify any existing .proto files from degrees
- ADD new .proto files for detailing domain services
- Always include google.api.http annotations for REST gateway mapping
- Use int64 for all IDs
- Use int64 for all money fields (cents, never floats)
- Follow existing degrees proto patterns for request/response naming
- Group related RPCs into a single service per domain
- All timestamps use google.protobuf.Timestamp
