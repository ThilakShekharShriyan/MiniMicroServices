package nodemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	NodeCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hashring_node_count",
		Help: "Number of nodes currently in the hash ring",
	})

	KeyLookups = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hashring_key_lookups_total",
			Help: "Total number of key lookups by node",
		},
		[]string{"node"},
	)

	NodeAdditions = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hashring_node_additions_total",
		Help: "Total nodes added to the ring",
	})

	NodeRemovals = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hashring_node_removals_total",
		Help: "Total nodes removed from the ring",
	})
)

func InitMetrics() {
	prometheus.MustRegister(NodeCount, KeyLookups, NodeAdditions, NodeRemovals)
}
