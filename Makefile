
BINARY_PATH=build/yaac
SOURCE_PATH=cmd/yaac

# The commented out line below should work on systems that use GNU sed
# CGO_CPPFLAGS=$(foreach dir,$(shell pkg-config --cflags lept tesseract), $(shell echo "$(dir)" | sed -E 's|(([a-zA-Z]+)\/[0-9.]+\/include\/)\2\/*|\1|g'))
CGO_LDFLAGS=$(shell pkg-config --libs lept tesseract)
# The line below is a hacked solution, it works for now,
# but should be removed in the future for a sed solution as above, that is portable
CGO_CPPFLAGS=$(foreach dir,$(shell pkg-config --cflags lept tesseract), $(patsubst %include/leptonica, %include, $(dir)))

.PHONY: all build test run clean

yaac: $(SOURCE_PATH)/*.go
	go build -o ./$(BINARY_PATH) ./$(SOURCE_PATH)

all: build test

build:
	go build -o ./$(BINARY_PATH) ./$(SOURCE_PATH)

test:
	go test -v ./test/
	
build-macos:
	CGO_CPPFLAGS=$(CGO_CPPFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go build -o ./$(BINARY_PATH) ./$(SOURCE_PATH)
	
test-macos:
	CGO_CPPFLAGS=$(CGO_CPPFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go test -v ./test/

run: build
	./$(BINARY_PATH)

package:
	fyne package --src ./$(SOURCE_PATH)/ --icon ./../../Icon.png

clean:
	go clean
	rm ./$(BINARY_PATH)

build-docker:
	docker build -t yaac-image .

run-docker:
	docker run -it --rm -e DISPLAY=$$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix yaac-image
