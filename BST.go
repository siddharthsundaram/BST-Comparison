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
var trees [] *TreeNode							// Slice of TreeNode ptrs
var hash_map map[int] []int						// Map of int -> slice of ints
var identical_trees [][] bool
var identical_tree_groups [] map[int] struct{}	// Slice of sets (implemented thru maps)
var map_lock sync.Mutex							// Mutex for hash map synchronization
var proceed sync.Cond

type Pair struct {
	Hash int
	Idx int
} 

func thread_hash(idx int, hashes []int, finished_barrier *sync.WaitGroup, hash_barrier *sync.WaitGroup, ch chan Pair, helper bool) {
	defer finished_barrier.Done()

	for i := idx; i < len(trees); i += *hash_workers {
		hash := 1
		trees[i].ComputeHash(&hash)
		if helper == true {

			// Send hash and index to central manager goroutine for grouping
			ch <- Pair{hash, i}
		} else {

			// Hash has been computed, store in hashes array
			hashes[i] = hash
		}
	}

	if helper == false {

		// Inform main thread that all hashes have been computed and block until hash computation time is recorded
		hash_barrier.Done()
		proceed.L.Lock()
		proceed.Wait()
		proceed.L.Unlock()

		// Time has been recorded, goroutines have been broadcasted to and will now group hashes
		for i := idx; i < len(hashes); i += *hash_workers {

			// Hash grouping begins
			map_lock.Lock()
			_, ok := hash_map[hashes[i]]
			if ok == true {
				hash_map[hashes[i]] = append(hash_map[hashes[i]], i)
			} else {
				hash_map[hashes[i]] = []int{i}
			}

			map_lock.Unlock()
		}
	}
}

func central_hash_group(ch chan Pair, hash_barrier *sync.WaitGroup) {

	// Group hashes and signal main thread after completion
	defer hash_barrier.Done()
	for pair := range ch {
		hash := pair.Hash
		idx := pair.Idx

		_, ok := hash_map[hash]
		if ok == true {
			hash_map[hash] = append(hash_map[hash], idx)
		} else {
			hash_map[hash] = []int{idx}
		}
	}
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

func print_hash_groups() {
	for key, val := range hash_map {
		if len(val) > 1 {
			fmt.Print(key, ": ")
			for j := range val {
				fmt.Print(val[j], " ")
			}
			fmt.Println()
		}
	}
}

func print_identical_matrix() {
	for i := range identical_trees {
		for j := range identical_trees {
			fmt.Print(identical_trees[i][j], " ")
		}

		fmt.Println()
	}
}

func main() {
	ParseCLI()

	// Populate trees slice
	ReadFile()

	// Set up necessary data structures
	var proceed_lock sync.Mutex
	proceed.L = &proceed_lock
	hashes := make([]int, len(trees))
	var hash_barrier sync.WaitGroup
	var finished_barrier sync.WaitGroup
	ch := make(chan Pair)
	hash_map = make(map[int] []int)
	var start time.Time
	var hash_time, hash_group_time time.Duration
	n := len(trees)
	identical_trees = make([][]bool, n)
	for i:= range identical_trees {
		identical_trees[i] = make([]bool, n)
	}

	if *hash_workers == 1 && *data_workers == 1 {

		// Serial hash computation
		start = time.Now()
		serial_hash(hashes)
		hash_time = time.Since(start)

		// Serial hash grouping
		serial_hash_group(hashes)
		hash_group_time = time.Since(start)
		
	} else if *data_workers == 1 {

		// Parallel hash computation with central manager to group hashes
		start = time.Now()

		// Assuming that goroutine must be spawned to group hashes rather than the main thread doing it
		hash_barrier.Add(1)
		go central_hash_group(ch, &hash_barrier)

		finished_barrier.Add(*hash_workers)
		for i := 0; i < *hash_workers; i++ {
			go thread_hash(i, hashes, &finished_barrier, &hash_barrier, ch, true)
		}

		finished_barrier.Wait()
		hash_time = time.Since(start)
		close(ch)

		hash_barrier.Wait()
		hash_group_time = time.Since(start)

	} else if *hash_workers > 1 && *hash_workers == *data_workers {

		// Parallel hash computation where each goroutine adds hash to map
		start = time.Now()

		hash_barrier.Add(*hash_workers)
		finished_barrier.Add(*hash_workers)
		for i := 0; i < *hash_workers; i++ {
			go thread_hash(i, hashes, &finished_barrier, &hash_barrier, ch, false)
		}

		hash_barrier.Wait()
		hash_time = time.Since(start)
		// proceed.L.Lock()					// Not sure if this locking is necessary, ask in OH
		proceed.Broadcast()
		// proceed.Unlock()
		finished_barrier.Wait()
		hash_group_time = time.Since(start)
	}

	fmt.Printf("hashTime: %.10f\n", hash_time.Seconds())
	fmt.Printf("hashGroupTime: %.10f\n", hash_group_time.Seconds())
	print_hash_groups()

	// Compare tree time
	start = time.Now()
	if *comp_workers > 0 {

		// TODO: Efficiently determine identical trees, including edge case
		// where there are multiple identical groups in the same hash group
		
		// Serial
		for _, val := range hash_map {
				// same := make(map[int] struct{})
				// same[trees[0]] = struct{}{}
			for i := 0; i < len(val); i++ {

				// Set equivalence to itself in adj matrix
				identical_trees[val[i]][val[i]] = true

				// One goroutine per comparison
				for j := i + 1; j < len(val); j++ {
					go CompareTrees(val[i], val[j])
				}
			}
		}

		print_identical_matrix()
	}

}