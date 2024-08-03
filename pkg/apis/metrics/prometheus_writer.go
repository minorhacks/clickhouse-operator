// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"fmt"
	"github.com/minorhacks/clickhouse-operator/pkg/metrics"
	"strconv"
	"time"

	log "github.com/golang/glog"
	// log "k8s.io/klog"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/minorhacks/clickhouse-operator/pkg/util"
)

const (
	namespace = "chi"
	subsystem = "clickhouse"
)

const (
	// writeMetricWaitTimeout specifies how long to wait for metric being accepted by prometheus writer
	writeMetricWaitTimeout = 10 * time.Second
)

// CHIPrometheusWriter specifies writer to prometheus
type CHIPrometheusWriter struct {
	out  chan<- prometheus.Metric
	chi  *WatchedCHI
	host *WatchedHost
}

// NewCHIPrometheusWriter creates new CHI prometheus writer
func NewCHIPrometheusWriter(
	out chan<- prometheus.Metric,
	chi *WatchedCHI,
	host *WatchedHost,
) *CHIPrometheusWriter {
	return &CHIPrometheusWriter{
		out:  out,
		chi:  chi,
		host: host,
	}
}

// WriteMetrics pushes set of prometheus.Metric objects created from the ClickHouse system data
// Expected data structure: metric, value, description, type (gauge|counter)
// TODO add namespace handling. It is just skipped for now
func (w *CHIPrometheusWriter) WriteMetrics(data [][]string) {
	for _, metric := range data {
		if len(metric) < 2 {
			continue
		}

		// Build metric from data row
		name := metric[0]
		value := metric[1]
		desc := metric[2]
		var metricType prometheus.ValueType
		if metric[3] == "counter" {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}
		w.writeSingleMetricToPrometheus(
			name, desc,
			metricType, value,
			nil, nil)
	}
}

// WriteTableSizes pushes set of prometheus.Metric objects created from the ClickHouse system data
// Expected data structure: database, table, partitions, parts, bytes, uncompressed_bytes, rows
// TODO add namespace handling. It is just skipped for now
func (w *CHIPrometheusWriter) WriteTableSizes(data [][]string) {
	for _, metric := range data {
		if len(metric) < 2 {
			continue
		}
		labelNames := []string{"database", "table", "active"}
		labelValues := []string{metric[0], metric[1], metric[2]}
		w.writeSingleMetricToPrometheus(
			"table_partitions", "Number of partitions of the table",
			prometheus.GaugeValue, metric[3],
			labelNames, labelValues)
		w.writeSingleMetricToPrometheus(
			"table_parts", "Number of parts of the table",
			prometheus.GaugeValue, metric[4],
			labelNames, labelValues)
		w.writeSingleMetricToPrometheus(
			"table_parts_bytes", "Table size in bytes",
			prometheus.GaugeValue, metric[5],
			labelNames, labelValues)
		w.writeSingleMetricToPrometheus(
			"table_parts_bytes_uncompressed", "Table size in bytes uncompressed",
			prometheus.GaugeValue, metric[6],
			labelNames, labelValues)
		w.writeSingleMetricToPrometheus(
			"table_parts_rows", "Number of rows in the table",
			prometheus.GaugeValue, metric[7],
			labelNames, labelValues)
	}
}

// WriteSystemParts pushes set of prometheus.Metric object related to system.parts
func (w *CHIPrometheusWriter) WriteSystemParts(data [][]string) {
	var diskDataBytes, memoryPrimaryKeyBytesAllocated int64
	var err error
	m := make([]int64, 2)
	for _, t := range data {
		m[0] = 0
		m[1] = 0
		for i, v := range t[len(t)-len(m):] {
			m[i], err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				log.V(2).Infof("Error parsing metrics value for chi_metric_DiskDataBytes, chi_metric_MemoryPrimaryKeyBytesAllocated: %s\n", v)
			}
		}
		diskDataBytes += m[0]
		memoryPrimaryKeyBytesAllocated += m[1]
	}
	w.WriteMetrics([][]string{
		{
			"metric.DiskDataBytes", fmt.Sprintf("%d", diskDataBytes),
			"Total data size for all ClickHouse tables", "gauge",
		},
		{
			"metric.MemoryPrimaryKeyBytesAllocated", fmt.Sprintf("%d", memoryPrimaryKeyBytesAllocated),
			"Memory size allocated for primary keys", "gauge",
		},
	})
}

// WriteSystemReplicas writes system replicas
func (w *CHIPrometheusWriter) WriteSystemReplicas(data [][]string) {
	for _, metric := range data {
		labelNames := []string{"database", "table"}
		labelValues := []string{metric[0], metric[1]}
		w.writeSingleMetricToPrometheus(
			"system_replicas_is_session_expired", "Number of expired Zookeeper sessions of the table",
			prometheus.GaugeValue, metric[2],
			labelNames, labelValues)
	}
}

// WriteMutations writes mutations
func (w *CHIPrometheusWriter) WriteMutations(data [][]string) {
	for _, metric := range data {
		labelNames := []string{"database", "table"}
		labelValues := []string{metric[0], metric[1]}
		w.writeSingleMetricToPrometheus(
			"table_mutations", "Number of active mutations for the table",
			prometheus.GaugeValue, metric[2],
			labelNames, labelValues)
		w.writeSingleMetricToPrometheus(
			"table_mutations_parts_to_do", "Number of data parts that need to be mutated for the mutation to finish",
			prometheus.GaugeValue, metric[3],
			labelNames, labelValues)
	}
}

// WriteSystemDisks writes system disks
func (w *CHIPrometheusWriter) WriteSystemDisks(data [][]string) {
	for _, metric := range data {
		labelNames := []string{"disk"}
		labelValues := []string{metric[0]}
		w.writeSingleMetricToPrometheus(
			"metric_DiskFreeBytes", "Free disk space available from system.disks",
			prometheus.GaugeValue, metric[1],
			labelNames, labelValues)
		w.writeSingleMetricToPrometheus(
			"metric_DiskTotalBytes", "Total disk space available from system.disks",
			prometheus.GaugeValue, metric[2],
			labelNames, labelValues)
	}
}

// WriteDetachedParts writes detached parts
func (w *CHIPrometheusWriter) WriteDetachedParts(data [][]string) {
	for _, metric := range data {
		labelNames := []string{"database", "table", "disk", "reason"}
		labelValues := []string{metric[1], metric[2], metric[3], metric[4]}
		w.writeSingleMetricToPrometheus(
			"metric_DetachedParts", "Count of currently detached parts from system.detached_parts",
			prometheus.GaugeValue, metric[0],
			labelNames, labelValues)
	}
}

// WriteErrorFetch writes error fetch
func (w *CHIPrometheusWriter) WriteErrorFetch(fetchType string) {
	labelNames := []string{"fetch_type"}
	labelValues := []string{fetchType}
	w.writeSingleMetricToPrometheus(
		"metric_fetch_errors", "status of fetching metrics from ClickHouse 1 - unsuccessful, 0 - successful",
		prometheus.GaugeValue, "1",
		labelNames, labelValues)
}

// WriteOKFetch writes successful fetch
func (w *CHIPrometheusWriter) WriteOKFetch(fetchType string) {
	labelNames := []string{"fetch_type"}
	labelValues := []string{fetchType}
	w.writeSingleMetricToPrometheus(
		"metric_fetch_errors", "status of fetching metrics from ClickHouse 1 - unsuccessful, 0 - successful",
		prometheus.GaugeValue, "0",
		labelNames, labelValues)
}

func (w *CHIPrometheusWriter) appendHostLabel(labels, values []string) ([]string, []string) {
	return append(labels, "hostname"), append(values, w.host.Hostname)
}

func (w *CHIPrometheusWriter) getMandatoryLabelsAndValues() (labelNames []string, labelValues []string) {
	// Prepare mandatory set of labels
	labelNames, labelValues = metrics.GetMandatoryLabelsAndValues(w.chi)
	// Append current host label
	labelNames, labelValues = w.appendHostLabel(labelNames, labelValues)

	return labelNames, labelValues
}

func (w *CHIPrometheusWriter) writeSingleMetricToPrometheus(
	name string,
	desc string,
	metricType prometheus.ValueType,
	value string,
	optionalLabels []string,
	optionalLabelValues []string,
) {
	// Prepare mandatory set of labels
	labelNames, labelValues := w.getMandatoryLabelsAndValues()
	// Append optional labels
	labelNames = append(labelNames, optionalLabels...)
	labelValues = append(labelValues, optionalLabelValues...)

	floatValue, _ := strconv.ParseFloat(value, 64)
	metric, err := prometheus.NewConstMetric(
		newMetricDescriptor(name, desc, labelNames),
		metricType,
		floatValue,
		labelValues...,
	)
	if err != nil {
		log.Warningf("Error creating metric: %s err: %s", name, err)
		return
	}
	// Send metric into channel
	select {
	case w.out <- metric:
	case <-time.After(writeMetricWaitTimeout):
		log.Warningf("Error sending metric to the channel: %s", name)
	}
}

// newMetricDescriptor creates a new prometheus.Desc object
func newMetricDescriptor(name, help string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, util.BuildPrometheusMetricName(name)),
		help,
		util.BuildPrometheusLabels(labels...),
		nil,
	)
}
