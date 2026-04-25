PROJECT_ROOT=$(shell git rev-parse --show-toplevel)

DOCKER_REPO?=vukovic96/pack-optimizer
DOCKER_IMAGE_TAG?=latest

GOA_VERSION=v3.26.0
GOA_CMD=go run goa.design/goa/v3/cmd/goa
DESIGN_PKG=github.com/aleksandarv/pack-optimizer/design
DESIGN_OUTPUT?=$(PROJECT_ROOT)

API_URL?=http://localhost:8080

all: test coverage build

build:
	go build -o ./bin/pack-optimizer ./cmd/optimizer; \

test: goenv fmt vet gotest

gotest:
	@files=$$(go list ./... | grep -v /gen | grep -v /design | grep -v /cmd | grep -v /test); \
	go test -v -race -timeout=30s $$files

integration-test:
	API_URL=$(API_URL) go test -v ./test/integration/...

fmt:
	go fmt ./...

vet:
	go vet ./...

goenv:
	@go version

coverage:
	@files=$$(go list ./... | grep -v /gen | grep -v /design | grep -v /cmd | grep -v /test); \
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

goa-install:
	go install goa.design/goa/v3/cmd/goa@$(GOA_VERSION)

generate: goa-install
	$(GOA_CMD) gen $(DESIGN_PKG) -o $(DESIGN_OUTPUT)

docker-build:
	docker build . -t $(DOCKER_REPO):$(DOCKER_IMAGE_TAG)

push: docker-build
	docker push $(DOCKER_REPO):$(DOCKER_IMAGE_TAG)