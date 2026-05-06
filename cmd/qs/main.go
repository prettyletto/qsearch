package main

import (
	"fmt"
	"os"

	"github.com/prettyletto/qseach/internal/app/search"
	"github.com/prettyletto/qseach/internal/dispatch"
	"github.com/prettyletto/qseach/internal/domain/provider"
	"github.com/prettyletto/qseach/internal/infra/browser"
	"github.com/prettyletto/qseach/internal/providers/google"
)

func main() {
	providers := []provider.Provider{
		google.New(),
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
