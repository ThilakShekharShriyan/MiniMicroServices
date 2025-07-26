package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/thilakshekharshriyan/m/bench"
	"github.com/thilakshekharshriyan/m/kv"
)

func main() {
    // Common flags
    numKeys := flag.Int("n", 1e6, "number of keys per store")
    concurrency := flag.Int("c", 4, "number of goroutines")
    flag.Parse()

    // List all store types and their constructors
    stores := []struct {
        Name    string
        Factory func() kv.KVStore
    }{
        {Name: "hash",     Factory: func() kv.KVStore { return kv.NewHashStore() }},
        {Name: "bptree",   Factory: func() kv.KVStore { return kv.NewBTreeStore() }},
        //{Name: "lsm",      Factory: func() kv.KVStore { return kv.NewLSMStore() }},
        //{Name: "skiplist", Factory: func() kv.KVStore { return kv.NewSkipListStore() }},
        //{Name: "trie",     Factory: func() kv.KVStore { return kv.NewTrieStore() }},
    }

    var results []bench.Result

    for _, s := range stores {
        fmt.Printf("→ Benchmarking %-8s store...\n", s.Name)
        store := s.Factory()

        cfg := bench.Config{
            Store:       store,
            StoreName:   s.Name,
            NumKeys:     *numKeys,
            Concurrency: *concurrency,
            Workload:    bench.Random, // could also flag this
            KeySize:     16,            // bytes per key (adjust as needed)
            ValueSize:   128,           // bytes per value (adjust as needed)
        }

        res, err := bench.RunBenchmark(cfg)
        if err != nil {
            log.Fatalf("benchmark %s failed: %v", s.Name, err)
        }

        fmt.Printf("   %s: %.2f writes/sec, avg read latency %.2fms, mem alloc %d bytes\n",
            res.StoreName,
            res.WriteOpsPerSec,
            func() float64 {
                var sum time.Duration
                for _, d := range res.ReadLatencies {
                    sum += d
                }
                return float64(sum.Milliseconds()) / float64(len(res.ReadLatencies))
            }(),
            res.MemAllocBytes,
        )

        results = append(results, res)
    }

    // Write everything to bench/results.csv
    outPath := "bench/results.csv"
    if err := bench.WriteResults(outPath, results); err != nil {
        log.Fatalf("failed to write results: %v", err)
    }
    fmt.Printf("✅ All benchmarks complete; results written to %s\n", outPath)
}
