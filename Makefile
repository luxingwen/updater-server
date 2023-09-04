ROOT_DIR    = $(shell pwd)

.PHONY: server
client:
	go build -o updater-server main.go


.PHONY: build
build: server