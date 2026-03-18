SWAG := $(shell go env GOPATH)/bin/swag

.PHONY: docs build run

docs:
	$(SWAG) init \
		--generalInfo cmd/main.go \
		--dir . \
		--output docs \
		--parseDependency \
		--parseInternal

build: docs
	go build -o bin/aturdana ./cmd/main.go

run: docs
	go run ./cmd/main.go
