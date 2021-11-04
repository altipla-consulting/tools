
FILES = $(shell find . -type f -name '*.go' -not -path './tmp/*')

build:

lint:
	@go vet ./...

test:
	@go test ./...

update-deps:
	go get -u $(GO_GLOB)
	go mod download
	go mod tidy
	go test $(GO_GLOB)

gofmt:
	@gofmt -s -w $(FILES)
	@gofmt -r '&α{} -> new(α)' -w $(FILES)
	@impsort cmd -p tools.altipla.consulting
