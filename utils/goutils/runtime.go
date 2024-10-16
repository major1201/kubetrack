package goutils

import (
	"fmt"
	"runtime"
	"time"
)

var (
	startTime = time.Now()
	sysStatus = new(SysStatus)
)

// SysStatus describes system runtime environment
type SysStatus struct {
	Uptime       string
	NumGoroutine int

	// General statistics.
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

func updateSystemStatus() {
	sysStatus.Uptime = HumanReadableDuration(time.Now().Unix() - startTime.Unix())

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus.NumGoroutine = runtime.NumGoroutine()

	sysStatus.MemAllocated = FileSize(m.Alloc)
	sysStatus.MemTotal = FileSize(m.TotalAlloc)
	sysStatus.MemSys = FileSize(m.Sys)
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = FileSize(m.HeapAlloc)
	sysStatus.HeapSys = FileSize(m.HeapSys)
	sysStatus.HeapIdle = FileSize(m.HeapIdle)
	sysStatus.HeapInuse = FileSize(m.HeapInuse)
	sysStatus.HeapReleased = FileSize(m.HeapReleased)
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = FileSize(m.StackInuse)
	sysStatus.StackSys = FileSize(m.StackSys)
	sysStatus.MSpanInuse = FileSize(m.MSpanInuse)
	sysStatus.MSpanSys = FileSize(m.MSpanSys)
	sysStatus.MCacheInuse = FileSize(m.MCacheInuse)
	sysStatus.MCacheSys = FileSize(m.MCacheSys)
	sysStatus.BuckHashSys = FileSize(m.BuckHashSys)
	sysStatus.GCSys = FileSize(m.GCSys)
	sysStatus.OtherSys = FileSize(m.OtherSys)

	sysStatus.NextGC = FileSize(m.NextGC)
	sysStatus.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC
}

// GetSysStatus returns runtime information
func GetSysStatus() *SysStatus {
	updateSystemStatus()
	return sysStatus
}
