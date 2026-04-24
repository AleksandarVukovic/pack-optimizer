PROJECT_ROOT=$(shell git rev-parse --show-toplevel)

all: test coverage build

build:
	go build -o $(PROJECT_ROOT)/bin/pack-optimizer ./cmd/optimizer

test: goenv fmt vet test

gotest: 
	@files=$$(go list ./...); \
	go test -v -race -timeout=30s $$files

fmt:
	go fmt ./...

vet:
	go vet ./...

goenv:
	@go version

coverage:
	@files=$$(go list ./...); \
	go test -coverprofile=coverage.out $$files; \
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}'); \
	echo "Total coverage: $$total"; \
	if [ $$(echo "$$total < 70.0" | sed 's/%//g' | bc) -eq 1 ]; then \
		echo "ERROR: Coverage is below 70%!"; \
		exit 1; \
	fi

clean:
	rm -rf $(PROJECT_ROOT)/bin
	rm -rf coverage.out coverage.html

GOA_VERSION=v3.26.0
GOA_CMD=go run goa.design/goa/v3/cmd/goa
DESIGN_PKG=github.com/aleksandarv/pack-optimizer/design
DESIGN_OUTPUT?=$(PROJECT_ROOT)

goa-install:
	go install goa.design/goa/v3/cmd/goa@$(GOA_VERSION)

generate: goa-install
	$(GOA_CMD) gen $(DESIGN_PKG) -o $(DESIGN_OUTPUT)