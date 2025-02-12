import matplotlib.pyplot as plt
import numpy as np

def read_file(input_file):
    res = []
    with open(input_file, "r") as file:
        for line in file:
            res.append(float(line.strip()))

    return res

def compare_hash_computation(per_tree_file, hash_workers_file, sequential_file, n, test):
    per_tree_mean = np.mean(read_file(per_tree_file))
    hash_worker_times = read_file(hash_workers_file)
    sequential_mean = np.mean(read_file(sequential_file))

    per_tree_speedup = sequential_mean / per_tree_mean
    hash_worker_speedup = [sequential_mean / time for time in hash_worker_times]
    hash_workers = [2 ** i for i in range(1, n.bit_length())]

    plt.figure()
    plt.axhline(y=per_tree_speedup, color='r', linestyle='-', label="goroutine per BST")
    plt.plot(hash_workers, hash_worker_speedup, marker='o', linestyle='-', label="hash-workers goroutines")

    plt.xlabel("hash-workers")
    plt.ylabel("Speedup (T_serial / T_parallel)")
    plt.xscale("log", base=2)
    plt.xticks(hash_workers, hash_workers, rotation=45)
    plt.legend()
    plt.grid(True)
    plt.title("Parallel Hash Computation Speedup Comparison for " + test + ".txt")
    plt.savefig("performance_visualizations/hash_computation_performance/hash_computation_comparison_" + test)
    plt.close()

def compare_hash_grouping(channel_file, lock_file, sequential_file, n, test):
    channel_times = read_file(channel_file)
    lock_times = read_file(lock_file)
    sequential_mean = np.mean(read_file(sequential_file))

    channel_speedup = [sequential_mean / time for time in channel_times]
    lock_speedup = [sequential_mean / time for time in lock_times]
    hash_workers = [2 ** i for i in range(1, n.bit_length())]

    plt.figure()
    plt.plot(hash_workers, channel_speedup, color='r', marker='o', linestyle='-', label="data-workers=1 (central manager)")
    plt.plot(hash_workers, lock_speedup, marker='o', linestyle='-', label="data-workers=hash-workers")

    plt.xlabel("hash-workers")
    plt.ylabel("Speedup (T_serial / T_parallel)")
    plt.xscale("log", base=2)
    plt.xticks(hash_workers, hash_workers, rotation=45)
    plt.legend()
    plt.grid(True)
    plt.title("Parallel Hash Grouping Speedup Comparison for " + test + ".txt")
    plt.savefig("performance_visualizations/hash_grouping_performance/hash_grouping_comparison_" + test)
    plt.close()

def compare_tree_comparison(per_tree_file, buffer_file, sequential_file, n, test):
    per_tree_mean = np.mean(read_file(per_tree_file))
    buffer_times = read_file(buffer_file)
    sequential_mean = np.mean(read_file(sequential_file))

    per_tree_speedup = sequential_mean / per_tree_mean
    buffer_speedup = [sequential_mean / time for time in buffer_times]
    comp_workers = [2 ** i for i in range(1, n.bit_length())]

    plt.figure()
    plt.axhline(y=per_tree_speedup, color='r', linestyle='-', label="goroutine per tree pair")
    plt.plot(comp_workers, buffer_speedup, marker='o', linestyle='-', label="comp-workers goroutines (buffer)")

    plt.xlabel("comp-workers")
    plt.ylabel("Speedup (T_serial / T_parallel)")
    plt.xscale("log", base=2)
    plt.xticks(comp_workers, comp_workers, rotation=45)
    plt.legend()
    plt.grid(True)
    plt.title("Parallel Tree Comparison Speedup Comparison for " + test + ".txt")
    plt.savefig("performance_visualizations/tree_comparison_performance/tree_comparison_speedup_" + test)
    plt.close()


# Hash computation speedup graphs
per_tree_coarse = "performance_visualizations/hash_computation_performance/per_tree_parallel_hash_computation_coarse.txt"
hash_workers_coarse = "performance_visualizations/hash_computation_performance/hash_workers_parallel_hash_computation_coarse.txt"
sequential_coarse = "performance_visualizations/hash_computation_performance/sequential_hash_computation_coarse.txt"
compare_hash_computation(per_tree_coarse, hash_workers_coarse, sequential_coarse, 100, "coarse")

per_tree_fine = "performance_visualizations/hash_computation_performance/per_tree_parallel_hash_computation_fine.txt"
hash_workers_fine = "performance_visualizations/hash_computation_performance/hash_workers_parallel_hash_computation_fine.txt"
sequential_fine = "performance_visualizations/hash_computation_performance/sequential_hash_computation_fine.txt"
compare_hash_computation(per_tree_fine, hash_workers_fine, sequential_fine, 100000, "fine")

# Hash grouping speedup graphs
channel_coarse = "performance_visualizations/hash_grouping_performance/channel_parallel_hash_grouping_coarse.txt"
lock_coarse = "performance_visualizations/hash_grouping_performance/lock_parallel_hash_grouping_coarse.txt"
sequential_coarse = "performance_visualizations/hash_grouping_performance/sequential_hash_grouping_coarse.txt"
compare_hash_grouping(channel_coarse, lock_coarse, sequential_coarse, 100, "coarse")

channel_fine = "performance_visualizations/hash_grouping_performance/channel_parallel_hash_grouping_fine.txt"
lock_fine = "performance_visualizations/hash_grouping_performance/lock_parallel_hash_grouping_fine.txt"
sequential_fine = "performance_visualizations/hash_grouping_performance/sequential_hash_grouping_fine.txt"
compare_hash_grouping(channel_fine, lock_fine, sequential_fine, 100000, "fine")

# Tree comparison speedup graphs
per_tree_coarse = "performance_visualizations/tree_comparison_performance/per_tree_coarse.txt"
buffer_coarse = "performance_visualizations/tree_comparison_performance/buffer_coarse.txt"
sequential_coarse = "performance_visualizations/tree_comparison_performance/sequential_tree_comparison_coarse.txt"
compare_tree_comparison(per_tree_coarse, buffer_coarse, sequential_coarse, 100, "coarse")

per_tree_fine = "performance_visualizations/tree_comparison_performance/per_tree_fine.txt"
buffer_fine = "performance_visualizations/tree_comparison_performance/buffer_fine.txt"
sequential_fine = "performance_visualizations/tree_comparison_performance/sequential_tree_comparison_fine.txt"
compare_tree_comparison(per_tree_fine, buffer_fine, sequential_fine, 100000, "fine")
