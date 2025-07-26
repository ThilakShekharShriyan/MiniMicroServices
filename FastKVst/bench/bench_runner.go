package bench

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
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
    StoreName   string       // e.g. "Hash", "BPTree", ...
    NumKeys     int
    Concurrency int
    Workload    WorkloadType
    KeySize     int          // bytes per key
    ValueSize   int          // bytes per value
}

// Result holds measured metrics.
type Result struct {
    StoreName      string
    NumKeys        int
    Concurrency    int
    WriteOpsPerSec float64
    ReadLatencies  []time.Duration
    MemAllocBytes  uint64
}

// generateWorkload builds a slice of keys according to cfg.Workload.
func generateWorkload(cfg Config) []string {
    keys := make([]string, cfg.NumKeys)
    switch cfg.Workload {
    case Sequential:
        width := int(math.Log10(float64(cfg.NumKeys))) + 1
        for i := range keys {
            keys[i] = fmt.Sprintf("key-%0*d", width, i)
        }

    case Zipfian:
        // s=1.2, v=1, imax=NumKeys-1
        z := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 1.2, 1, uint64(cfg.NumKeys-1))
        for i := range keys {
            idx := z.Uint64()
            keys[i] = fmt.Sprintf("key-%06d", idx)
        }

    default: // Random
        letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
        for i := range keys {
            b := make([]rune, cfg.KeySize)
            for j := range b {
                b[j] = letters[rand.Intn(len(letters))]
            }
            keys[i] = string(b)
        }
    }
    return keys
}

// RunBenchmark runs concurrent Sets, then measures per‑Get latency.
func RunBenchmark(cfg Config) (Result, error) {
    // 1. Generate keys & a constant value
    keys := generateWorkload(cfg)
    valBytes := make([]byte, cfg.ValueSize)
    rand.Read(valBytes)
    val := string(valBytes)

    // 2. Concurrent writes
    startW := time.Now()
    var wg sync.WaitGroup
    wg.Add(cfg.Concurrency)
    for c := 0; c < cfg.Concurrency; c++ {
        go func(c int) {
            defer wg.Done()
            for i := c; i < cfg.NumKeys; i += cfg.Concurrency {
                if err := cfg.Store.Set(keys[i], val); err != nil {
                    panic(err) // you could collect errors instead
                }
            }
        }(c)
    }
    wg.Wait()
    writeDur := time.Since(startW)
    writeOpsPerSec := float64(cfg.NumKeys) / writeDur.Seconds()

    // 3. Force GC & read memory usage
    debug.FreeOSMemory()
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    // 4. Concurrent reads measuring per‑key latency
    readLatencies := make([]time.Duration, cfg.NumKeys)
    wg.Add(cfg.Concurrency)
    for c := 0; c < cfg.Concurrency; c++ {
        go func(c int) {
            defer wg.Done()
            for i := c; i < cfg.NumKeys; i += cfg.Concurrency {
                t0 := time.Now()
                if _, err := cfg.Store.Get(keys[i]); err != nil {
                    panic(err)
                }
                readLatencies[i] = time.Since(t0)
            }
        }(c)
    }
    wg.Wait()

    return Result{
        StoreName:      cfg.StoreName,
        NumKeys:        cfg.NumKeys,
        Concurrency:    cfg.Concurrency,
        WriteOpsPerSec: writeOpsPerSec,
        ReadLatencies:  readLatencies,
        MemAllocBytes:  m.Alloc,
    }, nil
}

// WriteResults writes bench results to a CSV at path.
func WriteResults(path string, results []Result) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    w := csv.NewWriter(f)
    defer w.Flush()

    // header
    if err := w.Write([]string{
        "Store", "NumKeys", "Concurrency", "Writes/sec",
        "AvgReadLatency(ms)", "MemAllocBytes",
    }); err != nil {
        return err
    }

    for _, r := range results {
        // average read latency in ms
        var sum time.Duration
        for _, d := range r.ReadLatencies {
            sum += d
        }
        avgMs := float64(sum.Milliseconds()) / float64(len(r.ReadLatencies))

        record := []string{
            r.StoreName,
            fmt.Sprintf("%d", r.NumKeys),
            fmt.Sprintf("%d", r.Concurrency),
            fmt.Sprintf("%.2f", r.WriteOpsPerSec),
            fmt.Sprintf("%.2f", avgMs),
            fmt.Sprintf("%d", r.MemAllocBytes),
        }
        if err := w.Write(record); err != nil {
            return err
        }
    }

    return nil
}
