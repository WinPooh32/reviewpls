package navigator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/WinPooh32/reviewpls/internal/source"
	"github.com/emirpasic/gods/v2/stacks/arraystack"
)

var ErrNoProvider = errors.New("no language provider defined")

type filename = string

type LanguageProvider interface {
	// Identifiers lists all identifiers at [Position] p.
	// Column offset is ignored.
	Identifiers(ctx context.Context, p Position) ([]Symbol, error)
	// Definition returns identifier definition by node at [Position] p.
	Definition(ctx context.Context, ident string, p Position) (Symbol, error)
	// Referencess returns identifier references by node at [Position] p.
	References(ctx context.Context, ident string, p Position) ([]Position, error)
	// Close cleans up provider resources.
	Close() error
}

type Position struct {
	File   filename
	Line   uint
	Column uint
}

type Meta struct {
	Description string
	Data        any
}

type Symbol struct {
	Docstring string
	Ident     string
	Pos       Position
	Meta      map[string]Meta
}

type Navigator struct {
	p         Position
	text      string
	moves     *arraystack.Stack[Position]
	providers map[source.Format]LanguageProvider
}

func NewNavigator() *Navigator {
	return &Navigator{
		p: Position{
			File:   "",
			Line:   0,
			Column: 0,
		},
		text:      "",
		moves:     arraystack.New[Position](),
		providers: map[source.Format]LanguageProvider{},
	}
}

func (nav *Navigator) AddLanguageProvider(lang source.Format, provider LanguageProvider) {
	nav.providers[lang] = provider
}

func (nav *Navigator) Identifiers(ctx context.Context) (idents []Symbol, err error) {
	provider, ok := nav.getProvider(nav.p.File)
	if !ok {
		return nil, fmt.Errorf("file %q: %w", nav.p.File, ErrNoProvider)
	}

	idents, err = provider.Identifiers(ctx, nav.p)
	if err != nil {
		return nil, fmt.Errorf("get identifiers at the position %+v: %w", nav.p, err)
	}

	return idents, nil
}

func (nav *Navigator) Definition(ctx context.Context, ident string) (symb Symbol, err error) {
	provider, ok := nav.getProvider(nav.p.File)
	if !ok {
		return symb, fmt.Errorf("file %q: %w", nav.p.File, ErrNoProvider)
	}

	symb, err = provider.Definition(ctx, ident, nav.p)
	if err != nil {
		return symb, fmt.Errorf("get symbol %q definition at the position %+v: %w", ident, nav.p, err)
	}

	return symb, nil
}

func (nav *Navigator) References(ctx context.Context, ident string, p Position) (pp []Position, err error) {
	provider, ok := nav.getProvider(nav.p.File)
	if !ok {
		return nil, fmt.Errorf("file %q: %w", nav.p.File, ErrNoProvider)
	}

	pp, err = provider.References(ctx, ident, nav.p)
	if err != nil {
		return nil, fmt.Errorf("get symbol %q references at the position %+v: %w", ident, nav.p, err)
	}

	return pp, nil
}

// GoTo moves cursor to the [Position] p.
func (nav *Navigator) GoTo(p Position) (err error) {
	text, err := nav.readAt(p)
	if err != nil {
		return err
	}

	nav.moves.Push(nav.p)
	nav.text = text
	nav.p = p

	return nil
}

// MoveBack returns true if successfully moved to the previous position.
func (nav *Navigator) MoveBack(p Position) (ok bool, err error) {
	mp, mok := nav.moves.Pop()
	if !mok {
		return false, nil
	}

	text, err := nav.readAt(p)
	if err != nil {
		return false, err
	}

	nav.text = text
	nav.p = mp

	return true, nil
}

func (nav *Navigator) readAt(p Position) (text string, err error) {
	content, err := os.ReadFile(p.File)
	if err != nil {
		return "", fmt.Errorf("read file %q: %w", p.File, err)
	}

	text = string(content)
	lines := strings.Count(text, "\n")

	if lines == 0 {
		return "", fmt.Errorf("file %q is empty", p.File)
	}

	if lines < 0 {
		panic("unreachable")
	}

	if p.Line >= uint(lines) {
		return "", fmt.Errorf("requested line %d is out of the file lines: %d", p.Line, lines)
	}

	return text, nil
}

func (nav *Navigator) Close() error {
	var errs []error

	for _, tree := range nav.providers {
		if err := tree.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (nav *Navigator) getProvider(file filename) (provider LanguageProvider, ok bool) {
	provider, ok = nav.providers[source.DetectByExt(file)]
	return
}
