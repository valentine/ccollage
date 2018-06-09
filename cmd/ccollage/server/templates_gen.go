package server

import (
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"

	"github.com/shurcooL/vfsgen"
)

// Generate uses the vfsgen package to bundle template files in the compiled code
func Generate() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("ERROR: Template directory not found.")
	}

	var templateDir = http.Dir(filepath.Join(path.Dir(filename), "templates"))

	err := vfsgen.Generate(templateDir, vfsgen.Options{
		BuildTags:    "!dev",
		VariableName: "templates",
		PackageName:  "server",
		Filename:     filepath.Join(path.Dir(filename), "templates_vfsdata.go"),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
