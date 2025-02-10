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
    plt.savefig("performance_visualizations/hash_computation_performance/hash_computation_comparison_" + test)
    plt.close()

def compare_hash_grouping(channel_file, lock_file, sequential_file, step, test):
    channel_times = read_file(channel_file)
    lock_times = read_file(lock_file)
    sequential_mean = np.mean(read_file(sequential_file))

    channel_speedup = [sequential_mean / time for time in channel_times]
    lock_speedup = [sequential_mean / time for time in lock_times]
    hash_workers = list(range(step, step * 11, step))

    plt.figure()
    plt.plot(hash_workers, channel_speedup, color='r', marker='o', linestyle='-', label="data-workers=1 (central manager)")
    plt.plot(hash_workers, lock_speedup, marker='o', linestyle='-', label="data-workers=hash-workers")

    plt.xlabel("hash-workers")
    plt.ylabel("Speedup (T_serial / T_parallel)")
    plt.legend()
    plt.grid(True)
    plt.xticks(hash_workers)
    plt.title("Parallel Hash Grouping Speedup Comparison for " + test + ".txt")
    plt.savefig("performance_visualizations/hash_grouping_performance/hash_grouping_comparison_" + test)
    plt.close()


# Hash computation speedup graphs
per_tree_coarse = "performance_visualizations/hash_computation_performance/per_tree_parallel_hash_computation_coarse.txt"
hash_workers_coarse = "performance_visualizations/hash_computation_performance/hash_workers_parallel_hash_computation_coarse.txt"
sequential_coarse = "performance_visualizations/hash_computation_performance/sequential_hash_computation_coarse.txt"
compare_hash_computation(per_tree_coarse, hash_workers_coarse, sequential_coarse, 10, "coarse")

per_tree_fine = "performance_visualizations/hash_computation_performance/per_tree_parallel_hash_computation_fine.txt"
hash_workers_fine = "performance_visualizations/hash_computation_performance/hash_workers_parallel_hash_computation_fine.txt"
sequential_fine = "performance_visualizations/hash_computation_performance/sequential_hash_computation_fine.txt"
compare_hash_computation(per_tree_fine, hash_workers_fine, sequential_fine, 10000, "fine")

# Hash grouping speedup graphs
channel_coarse = "performance_visualizations/hash_grouping_performance/channel_parallel_hash_grouping_coarse.txt"
lock_coarse = "performance_visualizations/hash_grouping_performance/lock_parallel_hash_grouping_coarse.txt"
sequential_coarse = "performance_visualizations/hash_grouping_performance/sequential_hash_grouping_coarse.txt"
compare_hash_grouping(channel_coarse, lock_coarse, sequential_coarse, 10, "coarse")

channel_fine = "performance_visualizations/hash_grouping_performance/channel_parallel_hash_grouping_fine.txt"
lock_fine = "performance_visualizations/hash_grouping_performance/lock_parallel_hash_grouping_fine.txt"
sequential_fine = "performance_visualizations/hash_grouping_performance/sequential_hash_grouping_fine.txt"
compare_hash_grouping(channel_fine, lock_fine, sequential_fine, 10000, "fine")