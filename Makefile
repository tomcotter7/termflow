.DEFAULT_GOAL := run
.PHONY: run build install

run:
	go run cmd/termflow/main.go

build:
	go build -o bin/termflow ./cmd/termflow/main.go

install: build
	mkdir -p $(HOME)/.local/bin
	ln -sf $(PWD)/bin/termflow $(HOME)/.local/bin/termflow
