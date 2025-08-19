package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/WinPooh32/reviewpls/internal/navigator"
	"github.com/WinPooh32/reviewpls/internal/navigator/providers/golang"
	"github.com/WinPooh32/reviewpls/internal/source"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	workspaceFolders := []string{"."}

	nav := navigator.NewNavigator()
	nav.AddLanguageProvider(Go(ctx, workspaceFolders))

	fmt.Println(nav.GoTo(navigator.Position{
		File:   "main.go",
		Line:   20,
		Column: 0,
	}))
	fmt.Println(nav.Identifiers(ctx))
}

// Go gets golang provider.
func Go(ctx context.Context, workspaceFolders []string) (source.Format, navigator.LanguageProvider) {
	provider, err := golang.New(ctx, workspaceFolders)
	if err != nil {
		panic(err)
	}

	return source.Go, provider
}
