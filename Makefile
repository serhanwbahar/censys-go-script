BIN=./bin/censyscli
CMD=./cmd/cli

run: 
	go run ${CMD} google.com

# Builds the binary
build:
	go build -o ${BIN} ${CMD}

# Creates the docker image
create-image:
	docker build --tag censyscli .

# Run all tests
test:
	go test ./pkg/... -v --count=1

# Shorthand for go mod tidy
tidy:
	go mod tidy

# Clean the binaries
clean:
	go clean
	rm ./bin/*

.PHONY: run test build tidy clean create-image