package apispec

// Import embed for the go:embed directive below
import _ "embed"

// OpenAPISpec contains the embedded OpenAPI specification file in YAML format
//
//go:embed openapi.yaml
var OpenAPISpec []byte
