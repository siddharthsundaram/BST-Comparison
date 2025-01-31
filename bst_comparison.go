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
var hash_map map[int][]int

func thread_hash(idx int, hash int, hashes []int, barrier *sync.WaitGroup) {
	defer barrier.Done()
	trees[idx].ComputeHash(&hash)
	hashes[idx] = hash
}

func serial_hash(hashes []int) {
	for i := range trees {
		var hash int = 1
		trees[i].ComputeHash(&hash)
		hashes[i] = hash
	}
}

func serial_hash_group(hashes []int) {
	for i := range hashes {
		_, ok := hash_map[hashes[i]]
		if ok == true {
			hash_map[hashes[i]] = append(hash_map[hashes[i]], i)
		} else {
			hash_map[hashes[i]] = []int{i}
		}
	}
}

func main() {
	ParseCLI()

	// Populate trees slice
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
		serial_hash(hashes)
	}

	elapsed := time.Since(start)
	fmt.Printf("hashTime: %.10f\n", elapsed.Seconds())

	// for i := range hashes {
	// 	fmt.Println(hashes[i])
	// }

	// Hash group time
	start = time.Now()
	if *data_workers > 0 {
		hash_map = make(map[int][]int)
		serial_hash_group(hashes)
	}

	elapsed = time.Since(start)
	fmt.Printf("hashGroupTime: %.10f\n", elapsed.Seconds())

	// Compare tree time
	// start = time.Now()
	// if *comp_workers > 0 {

	// }

}