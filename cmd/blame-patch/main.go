package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/WinPooh32/reviewpls/internal/gitrepo"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rootDir := flag.String("root-dir", "", "git root dir")
	masterBranch := flag.String("master", "master", "master branch")
	featBranch := flag.String("feature", "", "feature branch")
	blameFile := flag.String("blame", "", "file path relarive to root-dir")

	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	repo, err := gitrepo.OpenRepository(ctx, *rootDir)
	checkErr(err)

	defer checkErr(repo.Close())

	diffCommits, err := repo.DiffCommits(ctx, *masterBranch, *featBranch)
	checkErr(err)

	blamePatch, err := repo.BlamePatch(ctx, *featBranch, diffCommits, *blameFile)
	checkErr(err)

	fmt.Println(blamePatch)
}
