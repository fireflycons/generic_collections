GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

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
