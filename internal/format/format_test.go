package format_test

import (
	"testing"

	"github.com/WinPooh32/reviewpls/internal/format"
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
		want  format.Format
	}{
		{"binary", binaryFiles, format.Binary},
		{"text", textFiles, format.Text},
		{"go", goFiles, format.Go},
		{"c", cFiles, format.C},
		{"cpp", cppFiles, format.Cpp},
		{"python", pythonFiles, format.Python},
		{"javascript", javascriptFiles, format.Javascript},
		{"json", jsonFiles, format.Json},
		{"html", htmlFiles, format.Html},
		{"rust", rustFiles, format.Rust},
		{"unknown", []string{".blender"}, format.Unknown},
	}

	for _, tt := range tests {
		for _, filename := range tt.names {
			t.Run(tt.name+"_"+filename, func(t *testing.T) {
				t.Parallel()

				got := format.DetectByExt(filename)
				assert.Equal(t, tt.want.String(), got.String())
			})
		}
	}
}
