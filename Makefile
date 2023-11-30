GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

build:
	go build -v ./...

lint:
	golangci-lint run

test:
	go test -v -cover -timeout=120s ./...

.PHONY: test/cover
test/cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out

benchtool:
	$(MAKE) -C .build/benchmark_processor

.PHONY: test/bench
test/bench: benchtool
	$(shell $(GOTEST)  -benchmem -run== -bench ^Benchmark ./... | ".build/benchmark_processor/benchmark-processor" > benchmarks.md)
