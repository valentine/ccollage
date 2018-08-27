// +build dev

// This templates file is used for hot reloading in dev mode
package server

import (
	"net/http"
	"path"
	"path/filepath"
	"runtime"
)

var _, filename, _, _ = runtime.Caller(0)

var templates = http.Dir(filepath.Join(path.Dir(filename), "templates"))
