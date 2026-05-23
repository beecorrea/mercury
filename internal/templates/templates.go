package templates

import _ "embed"

//go:embed index.html
var IndexHTML []byte

//go:embed 404.html
var NotFoundHTML []byte
