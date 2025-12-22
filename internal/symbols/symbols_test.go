package symbols_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/WinPooh32/reviewpls/internal/symbols"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtRange(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})

	tempFile := filepath.Join(dir, "test.go")

	testCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	x := 42
	y := x + 1
	return
}`

	wantSymbols := []symbols.Line{
		{N: 1, Symbols: []symbols.Symbol{{Column: 9, Name: "main"}}},
		{N: 5, Symbols: []symbols.Symbol{{Column: 6, Name: "main"}}},
		{N: 6, Symbols: []symbols.Symbol{{Column: 2, Name: "fmt"}, {Column: 6, Name: "Println"}}},
		{N: 7, Symbols: []symbols.Symbol{{Column: 2, Name: "x"}}},
		{N: 8, Symbols: []symbols.Symbol{{Column: 2, Name: "y"}, {Column: 7, Name: "x"}}},
	}

	err := writeTestFile(tempFile, testCode)
	require.NoError(t, err)

	t.Run("Test full range as string", func(t *testing.T) {
		t.Parallel()

		const wantString = `1: [{9 main}]
5: [{6 main}]
6: [{2 fmt} {6 Println}]
7: [{2 x}]
8: [{2 y} {7 x}]
`

		result, err := symbols.AtRange(tempFile, 0, -1)
		require.NoError(t, err)

		assert.Equal(t, wantString, result.String())
	})

	t.Run("Test full range", func(t *testing.T) {
		t.Parallel()

		result, err := symbols.AtRange(tempFile, 0, -1)
		require.NoError(t, err)

		assert.Equal(t, wantSymbols, result)
	})

	t.Run("Test specific range", func(t *testing.T) {
		t.Parallel()

		result, err := symbols.AtRange(tempFile, 1, 5)
		require.NoError(t, err)

		assert.Equal(t, wantSymbols[0:2], result)
	})

	t.Run("Test invalid file path", func(t *testing.T) {
		t.Parallel()

		result, err := symbols.AtRange("/non/existent/file.go", 0, -1)

		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Test empty file", func(t *testing.T) {
		t.Parallel()

		emptyFile := filepath.Join(t.TempDir(), "empty.go")
		err := writeTestFile(emptyFile, "")
		require.NoError(t, err)

		result, err := symbols.AtRange(emptyFile, 0, -1)
		require.NoError(t, err)

		assert.Empty(t, result)
	})
}

func writeTestFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
