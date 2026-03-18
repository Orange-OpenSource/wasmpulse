package collector

import (
	"log"
	"strconv"
	"sync"

	"github.com/Orange-OpenSource/wasmpulse/release/discovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/process"
)

type PidCollector struct {
	pids []discovery.WasmProcessInfo
	mu   sync.RWMutex // Mutex to protect the pids slice

	// CPU metrics
	cpuGauge         *prometheus.GaugeVec
	cpuSecondsUser   *prometheus.GaugeVec
	cpuSecondsSystem *prometheus.GaugeVec
	cpuSecondsTotal  *prometheus.GaugeVec
	threadCount      *prometheus.GaugeVec

	// Memory metrics
	rssGauge   *prometheus.GaugeVec
	vmsGauge   *prometheus.GaugeVec
	swapGauge  *prometheus.GaugeVec
	hwmGauge   *prometheus.GaugeVec
	stackGauge *prometheus.GaugeVec

	// I/O metrics
	openFDs    *prometheus.GaugeVec
	readBytes  *prometheus.GaugeVec
	writeBytes *prometheus.GaugeVec
}

func NewPidCollector() *PidCollector {
	prometheus_labels := []string{"wasm_file", "runtime", "pid"}

	return &PidCollector{
		pids: []discovery.WasmProcessInfo{},

		// CPU
		cpuGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_cpu_usage_percent",
				Help: "CPU usage of a process.",
			},
			prometheus_labels,
		),
		cpuSecondsUser: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_cpu_seconds_user",
				Help: "Total user CPU time spent for this process.",
			},
			prometheus_labels,
		),
		cpuSecondsSystem: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_cpu_seconds_system",
				Help: "Total system CPU time spent for this process.",
			},
			prometheus_labels,
		),
		cpuSecondsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_cpu_seconds_total",
				Help: "Total user+system CPU time spent for this process.",
			},
			prometheus_labels,
		),
		threadCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_threads_count",
				Help: "Number of OS threads created by this process.",
			},
			prometheus_labels,
		),

		// Memory
		rssGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_resident_memory_bytes",
				Help: "Resident Set Size of the process.",
			}, prometheus_labels),
		vmsGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_virtual_memory_bytes",
				Help: "Virtual Memory Size of the process (important for JIT reservation).",
			}, prometheus_labels),
		hwmGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_memory_high_water_mark_bytes",
				Help: "High water mark memory of the process.",
			}, prometheus_labels),
		stackGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "memory_stack_bytes",
				Help: "Stack size of the process.",
			}, prometheus_labels),
		swapGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_swap_memory_bytes",
				Help: "Swap memory used by the process.",
			}, prometheus_labels),

		// I/O
		openFDs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_open_fds",
				Help: "Number of open file descriptors.",
			}, prometheus_labels),
		readBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_disk_read_bytes_total",
				Help: "Total bytes read from disk.",
			}, prometheus_labels),
		writeBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_disk_write_bytes_total",
				Help: "Total bytes written to disk.",
			}, prometheus_labels),
	}
}

func (c *PidCollector) UpdatePids(newProcs []discovery.WasmProcessInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	existingPids := make(map[string]bool)
	for _, proc := range c.pids {
		existingPids[proc.PID] = true
	}

	for _, proc := range newProcs {
		if !existingPids[proc.PID] {
			log.Printf("New PID discovered: %s (%s - %s). Adding to monitor list.", proc.PID, proc.FileName, proc.RuntimeName)
			c.pids = append(c.pids, proc)
		}
	}
}

func (c *PidCollector) GetPidCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.pids)
}

func (c *PidCollector) Describe(ch chan<- *prometheus.Desc) {
	c.rssGauge.Describe(ch)
	c.vmsGauge.Describe(ch)
	c.swapGauge.Describe(ch)
	c.hwmGauge.Describe(ch)
	c.stackGauge.Describe(ch)
	c.cpuGauge.Describe(ch)
	c.cpuSecondsUser.Describe(ch)
	c.cpuSecondsSystem.Describe(ch)
	c.cpuSecondsTotal.Describe(ch)
	c.threadCount.Describe(ch)
	c.openFDs.Describe(ch)
	c.readBytes.Describe(ch)
	c.writeBytes.Describe(ch)
}

func (c *PidCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var activeProcs []discovery.WasmProcessInfo

	for _, procInfo := range c.pids {
		pidInt, err := strconv.ParseInt(procInfo.PID, 10, 32)
		if err != nil {
			log.Printf("Error parsing PID '%s': %v", pidInt, err)
			continue
		}

		proc, err := process.NewProcess(int32(pidInt))

		labels := prometheus.Labels{
			"wasm_file": procInfo.FileName,
			"runtime":   procInfo.RuntimeName,
			"pid":       procInfo.PID,
		}

		if err != nil {
			log.Printf("Process %s terminated.", procInfo.PID)
			c.deleteMetrics(labels)
			continue
		}

		activeProcs = append(activeProcs, procInfo)

		if memInfo, err := proc.MemoryInfo(); err == nil {
			c.rssGauge.With(labels).Set(float64(memInfo.RSS))
			c.vmsGauge.With(labels).Set(float64(memInfo.VMS))
			c.swapGauge.With(labels).Set(float64(memInfo.Swap))
			c.stackGauge.With(labels).Set(float64(memInfo.Stack))
			c.hwmGauge.With(labels).Set(float64(memInfo.HWM))
		}

		if cpuPercent, err := proc.CPUPercent(); err == nil {
			c.cpuGauge.With(labels).Set(cpuPercent)
		}
		if times, err := proc.Times(); err == nil {
			c.cpuSecondsUser.With(labels).Set(times.User)
			c.cpuSecondsSystem.With(labels).Set(times.System)
			c.cpuSecondsTotal.With(labels).Set(times.User + times.System)
		}
		if numThreads, err := proc.NumThreads(); err == nil {
			c.threadCount.With(labels).Set(float64(numThreads))
		}

		if ioStats, err := proc.IOCounters(); err == nil {
			c.readBytes.With(labels).Set(float64(ioStats.ReadBytes))
			c.writeBytes.With(labels).Set(float64(ioStats.WriteBytes))
		}
	}

	c.pids = activeProcs

	c.rssGauge.Collect(ch)
	c.vmsGauge.Collect(ch)
	c.hwmGauge.Collect(ch)
	c.stackGauge.Collect(ch)
	c.swapGauge.Collect(ch)
	c.cpuGauge.Collect(ch)
	c.cpuSecondsUser.Collect(ch)
	c.cpuSecondsSystem.Collect(ch)
	c.cpuSecondsTotal.Collect(ch)
	c.threadCount.Collect(ch)
	c.openFDs.Collect(ch)
	c.readBytes.Collect(ch)
	c.writeBytes.Collect(ch)
}

func (c *PidCollector) deleteMetrics(labels prometheus.Labels) {
	c.rssGauge.Delete(labels)
	c.vmsGauge.Delete(labels)
	c.hwmGauge.Delete(labels)
	c.swapGauge.Delete(labels)
	c.stackGauge.Delete(labels)
	c.cpuGauge.Delete(labels)
	c.cpuSecondsUser.Delete(labels)
	c.cpuSecondsSystem.Delete(labels)
	c.cpuSecondsTotal.Delete(labels)
	c.threadCount.Delete(labels)
	c.openFDs.Delete(labels)
	c.readBytes.Delete(labels)
	c.writeBytes.Delete(labels)
}
