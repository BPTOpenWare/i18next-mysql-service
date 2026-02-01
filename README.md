
# i18next MySQL Service 

Initial development in progress, do not use!
This service exposes a gRPC API with a gRPC-Gateway REST interface for storing and retrieving i18next resource files. The API is defined in [proto/resources/resources.proto](proto/resources/resources.proto) under the `ResourcesAPI` service.

## Overview
- Stores page-level i18next JSON resources identified by `application`, `component`, `page`.
- Supports CRUD via REST and gRPC.
- Returns raw JSON for page lookups and structured metadata for UUID-based operations.

## REST Endpoints
- GET: `/v1/resources/{application}/{component}/{page}`
	- Returns: JSON (`application/json`) body containing resource payload.
- POST: `/v1/resources:retrieve`
	- Body: `{ "uuid": "<resource-uuid>" }`
	- Returns: `PageResourceDetails` (includes `uuid`, `resource`, timestamps, tags).
- POST: `/v1/resources`
	- Body: `CreateResourceRequest` with identity and `resource` JSON.
	- Returns: `PageResourceDetails` for the newly created resource.
- PUT: `/v1/resources/{uuid}`
	- Body: `UpdateResourceRequest` with updated `resource` JSON and optional metadata.
	- Returns: updated `PageResourceDetails`.
- DELETE: `/v1/resources/{uuid}`
	- Returns: empty body (`204` or `200` with empty).

## gRPC Methods (ResourcesAPI)
- `GetResource(GetResourceRequest) returns (google.api.HttpBody)`
- `GetResourceByUUID(GetResourceByUUIDRequest) returns (PageResourceDetails)`
- `CreatePageResource(CreateResourceRequest) returns (PageResourceDetails)`
- `UpdatePageResource(UpdateResourceRequest) returns (PageResourceDetails)`
- `DeletePageResource(DeleteResourceRequest) returns (google.protobuf.Empty)`

## Message Shapes
- `GetResourceRequest`: `application`, `component`, `page` (strings)
- `GetResourceByUUIDRequest`: `uuid` (string)
- `PageResourceDetails`:
	- `uuid`: string
	- `resource`: object (stored as `google.protobuf.Struct`), standard i18next JSON
	- `creationUID`: string
	- `creationTimeStamp`: timestamp
	- `revisionUID`: string
	- `revisionTimeStamp`: timestamp
	- `tags`: array of strings
- `CreateResourceRequest`:
	- `application`, `component`, `page`: strings
	- `resource`: object (i18next JSON)
	- `creationUID`: string (optional)
	- `tags`: array of strings (optional)
- `UpdateResourceRequest`:
	- `uuid`: string
	- `resource`: object (i18next JSON)
	- `revisionUID`: string (optional)
	- `tags`: array of strings (optional)
- `DeleteResourceRequest`:
	- `uuid`: string

## Examples
Retrieve page JSON:

```bash
curl -s "http://localhost:${GatewayServerPort}/v1/resources/myapp/header/login" \
	-H "Accept: application/json"
```

Retrieve by UUID:

```bash
curl -s "http://localhost:${GatewayServerPort}/v1/resources:retrieve" \
	-H "Content-Type: application/json" \
	-d '{"uuid":"123e4567-e89b-12d3-a456-426614174000"}'
```

Create a resource:

```bash
curl -s "http://localhost:${GatewayServerPort}/v1/resources" \
	-H "Content-Type: application/json" \
	-d '{
		"application":"myapp",
		"component":"header",
		"page":"login",
		"resource": {"en":{"title":"Welcome"}},
		"creationUID":"user-123",
		"tags":["public","v1"]
	}'
```

Update a resource:

```bash
curl -s -X PUT "http://localhost:${GatewayServerPort}/v1/resources/123e4567-e89b-12d3-a456-426614174000" \
	-H "Content-Type: application/json" \
	-d '{
		"uuid":"123e4567-e89b-12d3-a456-426614174000",
		"resource": {"en":{"title":"Welcome Back"}},
		"revisionUID":"user-456",
		"tags":["public","v2"]
	}'
```

Delete a resource:

```bash
curl -s -X DELETE "http://localhost:${GatewayServerPort}/v1/resources/123e4567-e89b-12d3-a456-426614174000"
```

## Generate Stubs & OpenAPI

```bash
 protoc -I ./proto --go_out ./proto --go_opt paths=source_relative --go-grpc_out ./proto --go-grpc_opt paths=source_relative --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative --openapi_out=./Docs/ ./proto/resources/resources.proto 
```

## Run
- Configure ports via environment:
	- `GrpcServerPort` (gRPC server)
	- `GatewayServerPort` (HTTP gateway)
- Start the service: `go build` and run the binary; the gateway serves REST endpoints and proxies to gRPC.

