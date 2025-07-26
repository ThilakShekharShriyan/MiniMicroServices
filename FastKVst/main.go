package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/thilakshekharshriyan/m/bench"
	"github.com/thilakshekharshriyan/m/kv"
)

func main() {
    var (
        storeType   = flag.String("store", "hash", "one of: hash, bptree, lsm, skiplist, trie")
        numKeys     = flag.Int("n", 1e6, "number of keys")
        concurrency = flag.Int("c", 4, "goroutines")
    )
    flag.Parse()

    // choose store
    var store kv.KVStore
    switch *storeType {
    case "hash":
        store = kv.NewHashStore()
    // case "bptree": ...
    default:
        log.Fatalf("unknown store %q", *storeType)
    }

    cfg := bench.Config{
        Store:       store,
        NumKeys:     *numKeys,
        Concurrency: *concurrency,
        Workload:    bench.Random,
    }
    res, err := bench.RunBenchmark(cfg)
    if err != nil {
        log.Fatal(err)
    }
    // append to CSV
    if err := bench.WriteResults("bench/results.csv", []bench.Result{res}); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Done: %+v\n", res)
}
