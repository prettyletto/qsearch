BINARY := qs
MAIN := ./cmd/qs
PREFIX ?= $(shell go env GOPATH)/bin
INSTALL_PATH := $(PREFIX)/$(BINARY)

.PHONY: help all build run test install uninstall clean init init-force path

help:
	@echo "QSearch development commands"
	@echo
	@echo "  make build       Build ./$(BINARY)"
	@echo "  make run         Run the Google TUI from source"
	@echo "  make test        Run Go tests"
	@echo "  make install     Install $(BINARY) to $(INSTALL_PATH)"
	@echo "  make uninstall   Remove $(INSTALL_PATH)"
	@echo "  make init        Create the user providers.toml if missing"
	@echo "  make init-force  Overwrite providers.toml with defaults"
	@echo "  make path        Print the install path"
	@echo "  make clean       Remove local build output"

all: test build

build:
	go build -o $(BINARY) $(MAIN)

run:
	go run $(MAIN) g

test:
	go test ./...

install:
	go install $(MAIN)
	@echo "Installed: $(INSTALL_PATH)"
	@echo "Make sure this directory is in PATH: $(PREFIX)"

uninstall:
	rm -f $(INSTALL_PATH)

init:
	go run $(MAIN) init

init-force:
	go run $(MAIN) init --force

path:
	@echo $(INSTALL_PATH)

clean:
	rm -f $(BINARY)
