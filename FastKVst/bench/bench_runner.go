package bench

import (
	"encoding/csv"
	"fmt"
	"os"

	"time"

	"github.com/thilakshekharshriyan/m/kv"
)

// WorkloadType defines key distribution.
type WorkloadType int
const (
    Random WorkloadType = iota
    Sequential
    Zipfian
)

// Config holds bench parameters.
type Config struct {
    Store       kv.KVStore
    NumKeys     int
    Concurrency int
    Workload    WorkloadType
}

// Result holds measured metrics.
type Result struct {
    StoreName       string
    NumKeys         int
    Concurrency     int
    WriteOpsPerSec  float64
    ReadLatencies   []time.Duration
    MemAllocBytes   uint64
}

// RunBenchmark executes one scenario and returns a Result.
func RunBenchmark(cfg Config) (Result, error) {
    // 1. Pre‑generate keys/values per cfg.Workload
    // 2. Use a sync.WaitGroup + go routines to Set()
    // 3. Measure time to complete writes → ops/sec
    // 4. Trigger GC & runtime.ReadMemStats()
    // 5. Spawn read goroutines measuring latency
    // 6. Collect and return Result
    return Result{}, nil
}

// WriteResults writes a slice of results to CSV.
func WriteResults(path string, results []Result) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()

    // header
    w.Write([]string{"Store", "NumKeys", "Concurrency", "Writes/sec", "AvgReadLatency(ms)", "MemAllocBytes"})
    for _, r := range results {
        // compute avg latency
        var sum time.Duration
        for _, d := range r.ReadLatencies {
            sum += d
        }
        avg := float64(sum.Milliseconds()) / float64(len(r.ReadLatencies))
        w.Write([]string{
            r.StoreName,
            fmt.Sprintf("%d", r.NumKeys),
            fmt.Sprintf("%d", r.Concurrency),
            fmt.Sprintf("%.2f", r.WriteOpsPerSec),
            fmt.Sprintf("%.2f", avg),
            fmt.Sprintf("%d", r.MemAllocBytes),
        })
    }
    return nil
}
