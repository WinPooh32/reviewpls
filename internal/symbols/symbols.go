package symbols

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

// Filter node kinds that we're interested in
var interestedKinds = map[string]struct{}{
	"identifier":         {},
	"function_call":      {},
	"package_identifier": {},
	"type_identifier":    {},
	"field_identifier":   {},
}

type Line struct {
	N       int
	Symbols []Symbol
}

type Symbol struct {
	Column uint
	Name   string
}

func AtRange(sourcePath string, begin, end int) ([]Line, error) {
	scanAll := begin == 0 && end < 1
	untilEnd := begin > 0 && end < 1

	// Initialize the parser
	parser := tree_sitter.NewParser()
	defer parser.Close()

	// Set language to Go
	language := tree_sitter_go.Language()
	if err := parser.SetLanguage(tree_sitter.NewLanguage(language)); err != nil {
		return nil, fmt.Errorf("set tree sitter go language: %w", err)
	}

	sourceCode, err := os.ReadFile(filepath.Clean(sourcePath))
	if err != nil {
		return nil, fmt.Errorf("read source file: %w", err)
	}

	// Parse the source code
	tree := parser.Parse(sourceCode, nil)
	defer tree.Close()

	lines := strings.Count(string(sourceCode), "\n")
	lineSymbols := make([]Line, 0, lines)

	for i := range lines {
		k := i + 1
		switch {
		case scanAll:
		case untilEnd:
			if k < begin {
				continue
			}
		case k < begin || k > end:
			continue
		}

		symbols, err := atLine(tree, sourceCode, i)
		if err != nil {
			return nil, err
		}

		if len(symbols) == 0 {
			continue
		}

		lineSymbols = append(lineSymbols, Line{
			N:       k,
			Symbols: symbols,
		})
	}

	return lineSymbols, nil
}

// atLine returns an array of symbols at the given line in a Golang source file.
func atLine(tree *tree_sitter.Tree, source []byte, line int) ([]Symbol, error) {
	if tree == nil || tree.RootNode() == nil {
		return nil, fmt.Errorf("failed to parse source")
	}

	rootNode := tree.RootNode()
	symbols := []Symbol{}

	if line < 0 {
		return nil, fmt.Errorf("line must be a positive value")
	}

	targetLine := uint(line)

	walk(rootNode, targetLine, func(node *tree_sitter.Node) bool {
		startPos := node.StartPosition()

		if startPos.Row == targetLine {
			nodeKind := node.Kind()

			if _, found := interestedKinds[nodeKind]; found {
				symbols = append(symbols, Symbol{Column: node.StartPosition().Column, Name: node.Utf8Text(source)})
			}
		}

		return true
	})

	return symbols, nil
}

// walk the syntax tree recursively
func walk(node *tree_sitter.Node, targetLine uint, visitor func(*tree_sitter.Node) bool) {
	if visitor(node) {
		for i := uint(0); i < node.ChildCount(); i++ {
			walk(node.Child(i), targetLine, visitor)
		}
	}
}
