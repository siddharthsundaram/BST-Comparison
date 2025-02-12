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

	./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=tests/fine.txt
	# ./$(BINARY) -hash-workers=100 -data-workers=1 -comp-workers=2 -input=tests/fine.txt >> performance_visualizations/tree_comparison_performance/per_tree_fine.txt
	# ./$(BINARY) -hash-workers=1000 -data-workers=1000 -comp-workers=1000 -input=tests/fine.txt >> fine_out.txt

	# ./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=tests/coarse.txt
	# ./$(BINARY) -hash-workers=100 -data-workers=1 -comp-workers=2 -input=tests/coarse.txt >> performance_visualizations/tree_comparison_performance/per_tree_coarse.txt
	# ./$(BINARY) -hash-workers=2 -data-workers=2 -comp-workers=2 -input=tests/coarse.txt

	# ./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=2 -input=tests/simple.txt

time:
	bash -c 'num=2; while [[ $$num -le 100 ]]; do \
		./$(BINARY) -hash-workers=$$num -data-workers=$$num -comp-workers=$$num -input=tests/coarse.txt >> performance_visualizations/tree_comparison_performance/buffer_coarse.txt; \
		((num *= 2)); \
	done'