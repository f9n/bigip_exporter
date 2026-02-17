package collector

import (
	"strings"
	"time"

	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
)

// A PoolMemberCollector implements the prometheus.Collector.
type PoolMemberCollector struct {
	metrics                 map[string]poolMemberMetric
	bigip                   *f5.Device
	partitionsList          []string
	collectorScrapeStatus   *prometheus.GaugeVec
	collectorScrapeDuration *prometheus.SummaryVec
}

type poolMemberMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBPoolStatsInnerEntries) float64
	valueType prometheus.ValueType
}

// NewPoolMemberCollector returns a collector that collecting pool member statistics
func NewPoolMemberCollector(bigip *f5.Device, namespace string, partitionsList []string) (*PoolMemberCollector, error) {
	var (
		subsystem  = "pool_member"
		labelNames = []string{"partition", "pool", "member"}
	)
	return &PoolMemberCollector{
		metrics: map[string]poolMemberMetric{
			"serverside_curConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_cur_conns"),
					"serverside_cur_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_curConns.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"serverside_maxConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_max_conns"),
					"serverside_max_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_maxConns.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_totConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_tot_conns"),
					"serverside_tot_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_totConns.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_bytesIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bytes_in"),
					"serverside_bytes_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsIn.Value / 8)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_bytesOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bytes_out"),
					"serverside_bytes_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsOut.Value / 8)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_pktsIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_pkts_in"),
					"serverside_pkts_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_pktsIn.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_pktsOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_pkts_out"),
					"serverside_pkts_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_pktsOut.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"curSessions": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "cur_sessions"),
					"cur_sessions",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.CurSessions.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"totRequests": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "tot_requests"),
					"tot_requests",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.TotRequests.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"status_availabilityState": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "status_availability_state"),
					"status_availability_state",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					if entries.Status_availabilityState.Description == "available" {
						return 1
					}
					return 0
				},
				valueType: prometheus.GaugeValue,
			},
			"status_enabledState": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "status_enabled_state"),
					"status_enabled_state",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					if entries.Status_enabledState.Description == "enabled" {
						return 1
					}
					return 0
				},
				valueType: prometheus.GaugeValue,
			},
		},
		collectorScrapeStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "collector_scrape_status",
				Help:      "collector_scrape_status",
			},
			[]string{"collector"},
		),
		collectorScrapeDuration: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "collector_scrape_duration",
				Help:      "collector_scrape_duration",
			},
			[]string{"collector"},
		),
		bigip:          bigip,
		partitionsList: partitionsList,
	}, nil
}

// Collect collects metrics for BIG-IP pool members.
func (c *PoolMemberCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	memberCount := 0

	// First, get all pools
	err, allPoolStats := c.bigip.ShowAllPoolStats()
	if err != nil {
		c.collectorScrapeStatus.WithLabelValues("pool_member").Set(float64(0))
		logger.Warn("Failed to get pool list for member statistics", "error", err)
	} else {
		// For each pool, get its member statistics
		for poolKey := range allPoolStats.Entries {
			// Extract pool name from key
			// Key format: "https://localhost/mgmt/tm/ltm/pool/~Partition~PoolName/stats"
			keyParts := strings.Split(poolKey, "/")
			if len(keyParts) < 2 {
				logger.Warn("Pool key has insufficient parts", "key", poolKey)
				continue
			}

			poolPath := keyParts[len(keyParts)-2]
			poolPathParts := strings.Split(poolPath, "~")
			if len(poolPathParts) < 2 {
				logger.Warn("Pool path has insufficient parts", "poolPath", poolPath)
				continue
			}

			partition := poolPathParts[1]
			poolName := poolPathParts[len(poolPathParts)-1]

			// Filter by partition if specified
			if c.partitionsList != nil && !stringInSlice(partition, c.partitionsList) {
				logger.Debug("Skipping pool due to partition filter", "partition", partition, "pool", poolName)
				continue
			}

			// Get member statistics for this specific pool
			// Note: f5er library expects format like "/Partition/PoolName" which it converts to "~Partition~PoolName"
			poolFullPath := "/" + partition + "/" + poolName
			err, poolMemberStats := c.bigip.ShowPoolMembersStats(poolFullPath)
			if err != nil {
				logger.Warn("Failed to get member statistics for pool", "pool", poolFullPath, "error", err)
				continue
			}

			logger.Debug("Pool member stats entries count", "pool", poolFullPath, "count", len(poolMemberStats.Entries))

			// Process each member in the pool
			for key, memberStats := range poolMemberStats.Entries {
				logger.Debug("Processing pool member key", "key", key)

				// Key format: "https://localhost/mgmt/tm/ltm/pool/~Partition~PoolName/members/~Partition~MemberName:Port/stats"
				keyParts := strings.Split(key, "/")
				if len(keyParts) < 2 {
					logger.Warn("Member key has insufficient parts", "key", key)
					continue
				}

				// Extract member name from path
				memberPath := keyParts[len(keyParts)-2]
				memberPathParts := strings.Split(memberPath, "~")
				if len(memberPathParts) < 2 {
					logger.Warn("Member path has insufficient parts", "memberPath", memberPath)
					continue
				}
				memberName := memberPathParts[len(memberPathParts)-1]

				logger.Debug("Collecting pool member metrics", "partition", partition, "pool", poolName, "member", memberName)
				memberCount++

				labels := []string{partition, poolName, memberName}
				for _, metric := range c.metrics {
					ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, metric.extract(memberStats.NestedStats.Entries), labels...)
				}
			}
		}
		c.collectorScrapeStatus.WithLabelValues("pool_member").Set(float64(1))
		logger.Info("Successfully fetched statistics for pool members", "member_count", memberCount)
	}

	elapsed := time.Since(start)
	c.collectorScrapeDuration.WithLabelValues("pool_member").Observe(float64(elapsed.Seconds()))
	c.collectorScrapeStatus.Collect(ch)
	c.collectorScrapeDuration.Collect(ch)
	logger.Debug("Getting pool member statistics", "duration", elapsed)
}

// Describe describes the metrics exported from this collector.
func (c *PoolMemberCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
	c.collectorScrapeStatus.Describe(ch)
	c.collectorScrapeDuration.Describe(ch)
}
