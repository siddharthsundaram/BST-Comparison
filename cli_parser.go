package main

import "flag"
import "fmt"

func Parse_CLI() {
	hash_workers = flag.Int("hash-workers", 0, "Number of threads for computing BST hashes")
	data_workers = flag.Int("data-workers", 0, "Number of threads for adding hashes to the map")
	comp_workers = flag.Int("comp-workers", 0, "Number of threads for comparing BSTs with identical hashes")
	input_file = flag.String("input", "", "String path to input file")
	flag.Parse()
}

func Print_args() {
	fmt.Println("hash workers:", *hash_workers)
	fmt.Println("data workers:", *data_workers)
	fmt.Println("comp workers:", *comp_workers)
	fmt.Println("input file:", *input_file)
}