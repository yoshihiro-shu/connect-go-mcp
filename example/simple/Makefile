.PHONY: generate build run clean

# Generate code from proto files
generate:
	buf generate

# Build the example
build: generate
	go build -o bin/example .

# Run the example
run: build
	./bin/example

# Clean generated files and binaries
clean:
	rm -rf gen/
	rm -rf bin/

# Initialize dependencies
deps:
	go mod tidy

# Check if everything compiles
check: generate
	go build .