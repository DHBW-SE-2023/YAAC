BINARY_PATH=build/yaac
SOURCE_PATH=cmd/yaac

.PHONY: all build test run clean

yaac: $(SOURCE_PATH)/*.go
	go build -o ./$(BINARY_PATH) ./$(SOURCE_PATH)

all: build test

build:
	go build -o ./$(BINARY_PATH) ./$(SOURCE_PATH)

test:
	go test -v ./test/

run: build
	./$(BINARY_PATH)

package:
	fyne package --src ./$(SOURCE_PATH)/ --icon ./../../Icon.png

clean:
	go clean
	rm ./$(BINARY_PATH)