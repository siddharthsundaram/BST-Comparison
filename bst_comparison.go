package main

import (
	"fmt"
	"sync"
	"time"
)

var hash_workers *int
var data_workers *int
var comp_workers *int
var input_file *string
var trees [] *TreeNode
var hash_map map[int]int

func thread_hash(idx int, hash int, hashes []int, barrier *sync.WaitGroup) {
	defer barrier.Done()
	trees[idx].ComputeHash(&hash)
	hashes[idx] = hash
}

func main() {
	ParseCLI()
	// PrintArgs()

	if *data_workers > 0 {
		hash_map = make(map[int]int)
	}

	// Populates trees slice
	ReadFile()
	hashes := make([]int, len(trees))

	// Hash time
	start := time.Now()
	if *hash_workers > 1 {

		// Parallel
		var barrier sync.WaitGroup

		for i := range trees {
			barrier.Add(1)
			go thread_hash(i, 1, hashes, &barrier)
		}

		barrier.Wait()
	} else {

		// Serial
		for i := range trees {
			var hash int = 1
			trees[i].ComputeHash(&hash)
			hashes[i] = hash
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("hashTime: %.10f\n", elapsed.Seconds())

	// for i := range hashes {
	// 	fmt.Println(hashes[i])
	// }
}