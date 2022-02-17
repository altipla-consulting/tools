package embed

import _ "embed"

//go:embed jnet.jsonnet
var NativeFuncs []byte
