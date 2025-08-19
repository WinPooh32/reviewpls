package source_test

import (
	"testing"

	"github.com/WinPooh32/reviewpls/internal/source"
	"github.com/stretchr/testify/assert"
)

var (
	binaryFiles = []string{
		"main.exe",
		"/lib.dll",
		"/pkg/my-package/lib.so",
	}
	textFiles = []string{
		"main.csv",
		"/main.txt",
		"/pkg/my-package/source.txt",
	}
	goFiles = []string{
		"main.go",
		"/main.go",
		"/pkg/my-package/source.go",
	}
	cFiles = []string{
		"main.c",
		"main.h",
		"/main.c",
		"/main.h",
		"/pkg/my-package/source.c",
		"/pkg/my-package/source.h",
	}
	cppFiles = []string{
		"main.cpp",
		"main.hpp",
		"/main.cpp",
		"/main.hpp",
		"/pkg/my-package/source.cpp",
		"/pkg/my-package/source.hpp",
	}
	htmlFiles = []string{
		"main.html",
		"main.htm",
		"/main.html",
		"/main.htm",
		"/pkg/my-package/source.html",
		"/pkg/my-package/source.htm",
	}
	pythonFiles = []string{
		"main.py",
		"/main.py",
		"/pkg/my-package/source.py",
	}
	javascriptFiles = []string{
		"main.js",
		"/main.js",
		"/pkg/my-package/source.js",
	}
	jsonFiles = []string{
		"main.json",
		"/main.json",
		"/pkg/my-package/source.json",
	}
	rustFiles = []string{
		"main.rs",
		"/main.rs",
		"/pkg/my-package/source.rs",
	}
)

func TestDetectByExt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		names []string
		want  source.Format
	}{
		{"binary", binaryFiles, source.Binary},
		{"text", textFiles, source.Text},
		{"go", goFiles, source.Go},
		{"c", cFiles, source.C},
		{"cpp", cppFiles, source.Cpp},
		{"python", pythonFiles, source.Python},
		{"javascript", javascriptFiles, source.Javascript},
		{"json", jsonFiles, source.Json},
		{"html", htmlFiles, source.Html},
		{"rust", rustFiles, source.Rust},
		{"unknown", []string{".blender"}, source.Unknown},
	}

	for _, tt := range tests {
		for _, filename := range tt.names {
			t.Run(tt.name+"_"+filename, func(t *testing.T) {
				t.Parallel()

				got := source.DetectByExt(filename)
				assert.Equal(t, tt.want.String(), got.String())
			})
		}
	}
}
