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

build-docker:
	docker build -t yaac-image .

run-docker:
	docker run -it --rm -e DISPLAY=$$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix yaac-image

graph:
	graphpkg -stdout github.com/DHBW-SE-2023/YAAC > graph_all.svg
	graphpkg -stdout -match 'YAAC' github.com/DHBW-SE-2023/YAAC > graph_yaac.svg
	graphpkg -stdout -match 'YAAC|fyne|go-imap|sqlite|gosseract|gocv|gorm' github.com/DHBW-SE-2023/YAAC > graph_direct.svg