# Makefile

.PHONY: build clean test

build:
	go build -o steampipe-plugin-clevercloud main.go

clean:
	go clean
	rm -f steampipe-plugin-clevercloud

test:
	go test ./... -v

STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/tdannenmuller/clevercloud@latest/steampipe-plugin-clevercloud.plugin -tags "${BUILD_TAGS}" *.go
	steampipe service restart