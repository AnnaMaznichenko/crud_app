.PHONY: generate
generate:
	cd go && go generate ./...

.PHONY: test
test: generate
	cd go && go test ./...

.PHONY: deps
deps:
	cd go && go mod vendor

.PHONY: init
init: deps
	cp .env.example .env