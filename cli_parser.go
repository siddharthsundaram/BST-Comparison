package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func ParseCLI() {
	hash_workers = flag.Int("hash-workers", 1, "Number of threads for computing BST hashes")
	data_workers = flag.Int("data-workers", 0, "Number of threads for adding hashes to the map")
	comp_workers = flag.Int("comp-workers", 0, "Number of threads for comparing BSTs with identical hashes")
	input_file = flag.String("input", "", "String path to input file")
	flag.Parse()
}

func PrintArgs() {
	fmt.Println("hash workers:", *hash_workers)
	fmt.Println("data workers:", *data_workers)
	fmt.Println("comp workers:", *comp_workers)
	fmt.Println("input file:", *input_file)
}

func ReadFile() {
	file, err := os.Open(*input_file)
	if err != nil {
		fmt.Println("Error opening input file")
		return
	}

	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, " ")
		root_num, _ := strconv.Atoi(splits[0])
		root := NewTreeNode(root_num)
		trees = append(trees, root)
		var num int

		for i := 1; i < len(splits); i++ {
			num, _ = strconv.Atoi(splits[i])
			root.Insert(num)
		}
	}
	
	// Uncomment below code to print in order traversals of all trees

	// for i := range trees {
	// 	trees[i].PrintInOrderTraversal()
	// 	fmt.Println()
	// }
}