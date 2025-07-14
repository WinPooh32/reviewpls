package format

import (
	"path/filepath"
	"strings"
)

type Format int

const (
	Unknown Format = iota
	Binary
	Text
	Go
	C
	Cpp
	Html
	Python
	Javascript
	Json
	Rust
)

func (f Format) String() string {
	switch f {
	case Binary:
		return "binary"
	case Text:
		return "text"
	case Go:
		return "go"
	case C:
		return "c"
	case Cpp:
		return "cpp"
	case Html:
		return "html"
	case Python:
		return "python"
	case Javascript:
		return "javascript"
	case Json:
		return "json"
	case Rust:
		return "rust"
	case Unknown:
		fallthrough
	default:
		return "<unknown>"
	}
}

func DetectByExt(name string) Format {
	switch ext := strings.ToLower(filepath.Ext(name)); ext {
	case ".exe", ".dll", ".so", ".bin", ".png", ".jpg", ".jpeg":
		return Binary
	case ".go":
		return Go
	case ".c", ".h":
		return C
	case ".cpp", ".hpp", ".cxx", ".c++", ".cc", ".ixx":
		return Cpp
	case ".html", ".htm":
		return Html
	case ".py":
		return Python
	case ".js":
		return Javascript
	case ".json":
		return Json
	case ".rs":
		return Rust
	case ".txt", ".csv", ".yml", ".yaml", ".md", ".mod", ".sum":
		return Text
	}

	return Unknown
}
