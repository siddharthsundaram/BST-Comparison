import matplotlib.pyplot as plt
import numpy as np

def read_file(input_file):
    res = []
    with open(input_file, "r") as file:
        for line in file:
            res.append(float(line.strip()))

    return res

def compare_hash_computation(per_tree_file, hash_workers_file, sequential_file, step, test):
    per_tree_mean = np.mean(read_file(per_tree_file))
    hash_worker_times = read_file(hash_workers_file)
    sequential_mean = np.mean(read_file(sequential_file))

    per_tree_speedup = sequential_mean / per_tree_mean
    hash_worker_speedup = [sequential_mean / time for time in hash_worker_times]
    hash_workers = list(range(step, step * 11, step))

    plt.figure()
    plt.axhline(y=per_tree_speedup, color='r', linestyle='-', label="goroutine per BST")
    plt.plot(hash_workers, hash_worker_speedup, marker='o', linestyle='-', label="hash-workers goroutines")

    plt.xlabel("hash-workers")
    plt.ylabel("Speedup (T_serial / T_parallel)")
    plt.legend()
    plt.grid(True)
    plt.xticks(hash_workers)
    plt.title("Parallel Hash Computation Speedup Comparison for " + test + ".txt")
    plt.savefig("visualizations/hash_computation_comparison_" + test)
    plt.close()

per_tree_coarse = "visualizations/per_tree_parallel_hash_computation_coarse.txt"
hash_workers_coarse = "visualizations/hash_workers_parallel_hash_computation_coarse.txt"
sequential_coarse = "visualizations/sequential_hash_computation_coarse.txt"
compare_hash_computation(per_tree_coarse, hash_workers_coarse, sequential_coarse, 10, "coarse")

per_tree_fine = "visualizations/per_tree_parallel_hash_computation_fine.txt"
hash_workers_fine = "visualizations/hash_workers_parallel_hash_computation_fine.txt"
sequential_fine = "visualizations/sequential_hash_computation_fine.txt"
compare_hash_computation(per_tree_fine, hash_workers_fine, sequential_fine, 10000, "fine")