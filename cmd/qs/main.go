package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"qsearch/internal/provider/google"
)

func main() {
	query := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if query == "" {
		fmt.Fprintln(os.Stderr, "usage: qs search terms")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	suggestions, err := google.Suggestions(ctx, http.DefaultClient, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "google suggestions: %v\n", err)
		os.Exit(1)
	}

	if len(suggestions) == 0 {
		fmt.Println(query)
		return
	}

	for _, suggestion := range suggestions {
		fmt.Println(suggestion)
	}
}
