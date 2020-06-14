// +build tools

// buildtools package tracks dependencies to tools not referenced in any
// application code, yet actively used in the development process.
package buildtools

import (
	// used for OpenApi spec code generation
	_ "github.com/deepmap/oapi-codegen"
	// yaml cli processor similar to jq
	_ "github.com/mikefarah/yq/v3"
	// static file to GO embedding processor
	_ "github.com/markbates/pkger/cmd/pkger"
)
