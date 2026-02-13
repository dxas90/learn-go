package apispec

// _ imports embed package for the go:embed directive
import _ "embed"

// OpenAPISpec contains the embedded OpenAPI specification file in YAML format.
//
//go:embed openapi.yaml
var OpenAPISpec []byte
