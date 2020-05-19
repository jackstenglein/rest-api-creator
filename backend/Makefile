HANDLER_DIR := handlers
BIN_DIR := bin
HANDLER_FILES := $(wildcard $(HANDLER_DIR)/*.go)
BIN_FILES := $(patsubst $(HANDLER_DIR)/%.go, $(BIN_DIR)/%, $(HANDLER_FILES))

.PHONY: build clean deploy gomodgen

build: gomodgen $(BIN_FILES)

$(BIN_DIR)/%: $(HANDLER_DIR)/%.go
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o $@ $<

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
