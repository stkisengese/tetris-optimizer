# Makefile for tetris-optimizer

# Compiler
GO = go

# Project name
NAME = tetris-optimizer

# Source files
SRC = main.go

# Test files
TEST_FILES = $(wildcard *_test.go)

# Build flags
BUILD_FLAGS = 

.PHONY: all build test clean run

all: build test

build:
	@$(GO) build $(BUILD_FLAGS) -o $(NAME) $(SRC)

test:
	@$(GO) test -v $(TEST_FILES) $(SRC)

clean:
	@rm -f $(NAME)
	@rm -f *.out

run: build
	@./$(NAME)

# Format code
fmt:
	@$(GO) fmt ./...

# Run with sample file
sample: build
	@./$(NAME) sample.txt

# Install dependencies (if any)
# install:
#	@$(GO) mod download

# Lint check (if you want to add later)
# lint:
#	@golangci-lint run