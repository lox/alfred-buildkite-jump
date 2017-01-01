SOURCES := $(wildcard *.go)
BIN := alfred-buildkite-jump
FILES := $(BIN) info.plist icon.png
WORKFLOW := Buildkite\ Jump.alfredworkflow

$(WORKFLOW): $(FILES)
	zip -j "$@" $^

build: $(BIN)

$(BIN): $(SOURCES)
	go build -o $(BIN) $(SOURCES)

clean:
	rm $(BIN) $(WORKFLOW)