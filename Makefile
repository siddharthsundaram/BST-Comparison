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
	./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=fine.txt
	./$(BINARY) -hash-workers=100 -data-workers=1 -comp-workers=2 -input=fine.txt
	./$(BINARY) -hash-workers=100 -data-workers=100 -comp-workers=2 -input=fine.txt