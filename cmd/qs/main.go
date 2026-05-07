package main

import (
	"fmt"
	"os"

	"github.com/prettyletto/qsearch/internal/app/search"
	"github.com/prettyletto/qsearch/internal/config"
	"github.com/prettyletto/qsearch/internal/dispatch"
	"github.com/prettyletto/qsearch/internal/domain/provider"
	"github.com/prettyletto/qsearch/internal/infra/browser"
	"github.com/prettyletto/qsearch/internal/providers/google"
	"github.com/prettyletto/qsearch/internal/providers/youtube"
	"github.com/prettyletto/qsearch/internal/providers/ytmusic"
)

func main() {
	providers := []provider.Provider{
		google.New(),
		youtube.New(),
		ytmusic.New(),
	}

	customProviders, err := config.LoadProviders("configs/providers.example.toml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	providers = append(providers, customProviders...)

	opener := browser.NewOpener()
	searchRunner := search.NewRunner(opener, providers)

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
