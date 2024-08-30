package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	OperatingSystem string  `json:"operating_system"`
	Architecture    string  `json:"architecture"`
	Hostname        string  `json:"hostname"`
	NumberOfCPUs    int     `json:"number_of_cpus"`
	CPUModel        string  `json:"cpu_model"`
	TotalMemoryGB   float64 `json:"total_memory_gb"`
	FreeMemoryGB    float64 `json:"free_memory_gb"`
	UsedMemoryGB    float64 `json:"used_memory_gb"`
	Uptime          string  `json:"uptime"`
	BootTime        string  `json:"boot_time"`
	OS              string  `json:"os"`
	Platform        string  `json:"platform"`
	PlatformFamily  string  `json:"platform_family"`
	PlatformVersion string  `json:"platform_version"`
	KernelVersion   string  `json:"kernel_version"`
}

func getSystemInfo() (*SystemInfo, error) {
	// Get CPU info
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	// Get memory info
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Get host info
	hostStat, err := host.Info()
	if err != nil {
		return nil, err
	}

	sysInfo := &SystemInfo{
		OperatingSystem: runtime.GOOS,
		Architecture:    runtime.GOARCH,
		Hostname:        getHostname(),
		NumberOfCPUs:    runtime.NumCPU(),
		CPUModel:        cpuInfo[0].ModelName,
		TotalMemoryGB:   float64(vmStat.Total) / 1e9,
		FreeMemoryGB:    float64(vmStat.Free) / 1e9,
		UsedMemoryGB:    float64(vmStat.Used) / 1e9,
		Uptime:          (time.Duration(hostStat.Uptime) * time.Second).String(),
		BootTime:        time.Unix(int64(hostStat.BootTime), 0).String(),
		OS:              hostStat.OS,
		Platform:        hostStat.Platform,
		PlatformFamily:  hostStat.PlatformFamily,
		PlatformVersion: hostStat.PlatformVersion,
		KernelVersion:   hostStat.KernelVersion,
	}

	return sysInfo, nil
}

func systemInfoHandler(w http.ResponseWriter, r *http.Request) {
	sysInfo, err := getSystemInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sysInfo)
}

func main() {
	http.HandleFunc("/systeminfo", systemInfoHandler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}
	return hostname
}
