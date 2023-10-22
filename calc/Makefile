EXE_PATH = bin/calc
PACKAGE_TO_TEST = ./...
INTEGRATION_TEST_FILE=pkg/tests.sh
OUT_FILE = coverage.out
HTML = coverage.html

all: build test_with_coverage

build: $(ENTRYPOINT)
	go build -o $(EXE_PATH) cmd/main.go

run: build
	$(EXE_PATH)

clean:
	rm -rf bin/

rebuild: clean build

test:
	go test ./... -race -covermode=atomic -coverpkg $(PACKAGE_TO_TEST) -coverprofile=$(OUT_FILE)

test_with_coverage: test
	go tool cover -html $(OUT_FILE) -o $(HTML)

cleantest:
	rm coverage*

retest:
	make cleantest
	make
