BINARY = bst_comparison

SRC = $(wildcard *.go)

all: build

build:
	go build -o $(BINARY) $(SRC)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=simple.txt
	./$(BINARY) -hash-workers=2 -data-workers=2 -comp-workers=2 -input=simple.txt