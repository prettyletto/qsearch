package main

import (
	"fmt"
	"os"

	"github.com/prettyletto/qsearch/internal/app/search"
	"github.com/prettyletto/qsearch/internal/dispatch"
	"github.com/prettyletto/qsearch/internal/domain/provider"
	"github.com/prettyletto/qsearch/internal/infra/browser"
	"github.com/prettyletto/qsearch/internal/providers/google"
	"github.com/prettyletto/qsearch/internal/providers/youtube"
)

func main() {
	providers := []provider.Provider{
		google.New(),
		youtube.New(),
	}

	opener := browser.NewOpener()
	searchRunner := search.NewRunner(opener)

	dispatcher, err := dispatch.NewDispatcher(searchRunner, providers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := dispatcher.Dispatch(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
