package apispec

import _ "embed" // Required for go:embed directive below

// OpenAPISpec contains the embedded OpenAPI specification file in YAML format.
//
//go:embed openapi.yaml
var OpenAPISpec []byte
