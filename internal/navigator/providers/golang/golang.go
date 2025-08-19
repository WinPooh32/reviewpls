package golang

import (
	"context"
	"fmt"

	"github.com/WinPooh32/reviewpls/internal/navigator"
)

type Go struct {
	lsp *gopls
}

func New(ctx context.Context, workspaceFolders []string) (*Go, error) {
	lsp, err := newGopls(ctx, workspaceFolders)
	if err != nil {
		return nil, fmt.Errorf("new lsp server: %w", err)
	}

	return &Go{
		lsp: lsp,
	}, nil
}

func (g *Go) Identifiers(ctx context.Context, p navigator.Position) ([]navigator.Symbol, error) {
	const (
		Method        = 6
		Property      = 7
		Field         = 8
		Function      = 12
		Variable      = 13
		Constant      = 14
		TypeParameter = 26
	)

	allSymbols, err := g.lsp.documentSymbols(ctx, p.File)
	if err != nil {
		return nil, fmt.Errorf("get all document symbols: %w", err)
	}

	var lineSymbols []navigator.Symbol

	for _, symb := range allSymbols {
		switch symb.Kind {
		case Method,
			Property,
			Field,
			Function,
			Variable,
			Constant,
			TypeParameter:
			lineSymbols = append(lineSymbols, navigator.Symbol{
				Docstring: "",
				Ident:     symb.Name,
				Pos: navigator.Position{
					File:   p.File,
					Line:   p.Line,
					Column: symb.SelectionRange.Start.Character,
				},
				Meta: nil,
			})

		default:
			continue
		}
	}

	return lineSymbols, nil
}

func (g *Go) Definition(ctx context.Context, ident string, p navigator.Position) (navigator.Symbol, error) {
	panic("TODO: Implement")
}

func (g *Go) References(ctx context.Context, ident string, p navigator.Position) ([]navigator.Position, error) {
	panic("TODO: Implement")
}

func (g *Go) Close() error {
	panic("TODO: Implement")
}
