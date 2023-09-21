ifeq ($(OS),Windows_NT)
    target = benchmark-processor.exe
else
    target = benchmark-processor
endif

all: $(target)

$(target): benchmark_processor.go
	go build -o $(target) benchmark_processor.go
