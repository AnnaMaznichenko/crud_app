.PHONY: generate
generate:
	cd go && go generate ./...

.PHONY: test
test: generate
	cd go && go test ./...

.PHONY: init
init:
	cp .env.example .env