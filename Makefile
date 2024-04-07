# Binary name
BINARY_NAME=backup-restore-tool
MAIN_FILE=cmd/main.go

all: test build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)