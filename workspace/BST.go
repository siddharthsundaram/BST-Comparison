package main

import (
	"fmt"
	"sync"
	"time"
	"sort"
)

var hash_workers *int
var data_workers *int
var comp_workers *int
var input_file *string
var trees [] *TreeNode							// Slice of TreeNode ptrs
var hash_map map[int] []int						// Map of int -> slice of ints
var identical_trees [][] bool					// Adjacency matrix to store BST equivalency
var identical_tree_groups [] map[int] struct{}	// Slice of sets (implemented thru maps)
var map_lock sync.Mutex							// Mutex for hash map synchronization
var proceed sync.Cond							// Condition for precise timing
var buf *ConcurrentBuffer

// Used in parallelization to send hashing info through channel
type Pair struct {
	Hash int
	Idx int
} 

// Function for goroutines to compute BST hashes in parallel
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
		proceed.L.Lock()
		hash_barrier.Done()
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

// Groups hashes and signal main thread after completion
func central_hash_group(ch chan Pair, hash_barrier *sync.WaitGroup) {
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

// Function for goroutines to compare two BSTs from the buffer
func thread_compare(comp_barrier *sync.WaitGroup) {
	defer comp_barrier.Done()
	for {
		work, ok := buf.dequeue()
		if !ok {
			return
		}

		CompareTrees(work.idx1, work.idx2)
	}
}

// Wrapper for comparing BSTs
func per_tree_compare(comp_barrier *sync.WaitGroup, idx1 int, idx2 int) {
	defer comp_barrier.Done()
	CompareTrees(idx1, idx2)
}

// Serial hash computing wrapper function
func serial_hash(hashes []int) {
	for i := range trees {
		var hash int = 1
		trees[i].ComputeHash(&hash)
		hashes[i] = hash
	}
}

// Serial hash grouping function
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

// Function for printing the hash groups
func print_hash_groups() {
	for key, val := range hash_map {
		if len(val) > 1 {
			sort.Ints(val)
			fmt.Print(key, ": ")
			for j := range val {
				fmt.Print(val[j], " ")
			}
			fmt.Println()
		}
	}
}

// Function for printing equivalence adjacency matrix (for debugging)
func print_identical_matrix() {
	for i := range identical_trees {
		for j := range identical_trees {
			fmt.Print(identical_trees[i][j], " ")
		}

		fmt.Println()
	}
}

// Function for printing identical tree groups after comparison
func print_identical_groups() {
	seen := make(map[int]struct{})
	count := 0

	for i := range identical_trees {
		_, ok := seen[i]
		if ok {
			continue
		}

		group := [] int {i}
		seen[i] = struct{}{}
		for j := i + 1; j < len(identical_trees); j++ {
			if identical_trees[i][j] == true {
				group = append(group, j)
				seen[j] = struct{}{}
			}
		}

		if len(group) > 1 {
			fmt.Print("group", count, ": ")
			for j := range group {
				fmt.Print(group[j], " ")
			}

			count += 1
			fmt.Println()
		}

	}
}

func main() {
	e2e_start := time.Now()
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
	var hash_time, hash_group_time, comp_time time.Duration
	n := len(trees)
	identical_trees = make([][]bool, n)
	for i:= range identical_trees {
		identical_trees[i] = make([]bool, n)
	}

	// First implementation of parallelizing hash computation
	// *hash_workers = len(trees)
	// fmt.Println(*hash_workers)

	if *hash_workers == 1 {

		// Serial hash computation
		start = time.Now()
		serial_hash(hashes)
		hash_time = time.Since(start)

		// Serial hash grouping
		if *data_workers == 1 {
			serial_hash_group(hashes)
			hash_group_time = time.Since(start)
		}
		
	} else if *data_workers == 1 {

		// Parallel hash computation with central manager to group hashes
		start = time.Now()
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
		// proceed.L.Unlock()
		finished_barrier.Wait()
		hash_group_time = time.Since(start)
	}

	// fmt.Printf("hashTime: %.10f\n", hash_time.Seconds())
	// fmt.Printf("%.10f\n", hash_time.Seconds())
	// fmt.Printf("hashGroupTime: %.10f\n", hash_group_time.Seconds())
	// print_hash_groups()

	// Compare BSTs
	// Serial
	if *comp_workers == 1 {
		start = time.Now()

		for _, val := range hash_map {
			for i := 0; i < len(val); i++ {

				// Set equivalence to itself in adj matrix
				identical_trees[val[i]][val[i]] = true

				for j := i + 1; j < len(val); j++ {
					CompareTrees(val[i], val[j])
				}
			}
		}

		comp_time = time.Since(start)

	// Parallel
	} else if *comp_workers > 1 {
		start = time.Now()
		var comp_barrier sync.WaitGroup
		buf = NewConcurrentBuffer(*comp_workers)
		comp_barrier.Add(*comp_workers)
		for i := 0; i < *comp_workers; i++ {
			go thread_compare(&comp_barrier)
		}

		for _, val := range hash_map {
			for i := 0; i < len(val); i++ {

				// Set equivalence to itself in adj matrix
				identical_trees[val[i]][val[i]] = true

				for j := i + 1; j < len(val); j++ {

					// Implementation #1: One goroutine per comparison
					// comp_barrier.Add(1)
					// go per_tree_compare(&comp_barrier, val[i], val[j])

					// Implementation #2: Main thread puts work in concurrent buffer for goroutines to consume and perform
					buf.enqueue(val[i], val[j])
				}
			}
		}

		buf.Close()
		comp_barrier.Wait()
		comp_time = time.Since(start)
	}

	// fmt.Printf("compareTreeTime: %.10f\n", comp_time.Seconds())
	// fmt.Printf("%.10f\n", hash_time.Seconds())
	// fmt.Printf("%.10f\n", hash_group_time.Seconds())
	fmt.Printf("%.10f\n", comp_time.Seconds())
	// print_identical_groups()
	e2e_end := time.Since(e2e_start)
	// fmt.Printf("e2eTime: %.10f\n", e2e_end.Seconds())

	temp := hash_time
	temp = hash_group_time
	temp = comp_time
	temp = e2e_end
	comp_time = temp
}