BINARY = BST

SRC = $(wildcard workspace/*.go)

all: build

build:
	go build -o $(BINARY) $(SRC)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	# ./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=tests/simple.txt
	# ./$(BINARY) -hash-workers=5 -data-workers=1 -comp-workers=2 -input=tests/simple.txt
	# ./$(BINARY) -hash-workers=5 -data-workers=5 -comp-workers=2 -input=tests/simple.txt

	# ./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=tests/fine.txt
	# ./$(BINARY) -hash-workers=100 -data-workers=1 -comp-workers=2 -input=tests/fine.txt
	./$(BINARY) -hash-workers=1000 -data-workers=1000 -comp-workers=1000 -input=tests/fine.txt

	# ./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=tests/coarse.txt >> performance_visualizations/sequential_hash_grouping_coarse.txt
	# ./$(BINARY) -hash-workers=100 -data-workers=1 -comp-workers=2 -input=tests/coarse.txt
	# ./$(BINARY) -hash-workers=100 -data-workers=100 -comp-workers=2 -input=tests/coarse.txt

	# ./$(BINARY) -hash-workers=0 -data-workers=0 -comp-workers=0 -input=tests/simple.txt

time:
	bash -c 'num=10; while [[ $$num -le 100 ]]; do \
		./$(BINARY) -hash-workers=$$num -data-workers=$$num -comp-workers=2 -input=tests/coarse.txt >> performance_visualizations/lock_parallel_hash_grouping_coarse.txt; \
		((num += 10)); \
	done'