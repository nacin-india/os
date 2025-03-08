package system

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfo holds all system information
type SystemInfo struct {
	CPUInfo     string
	GPUInfo     string
	MemoryInfo  string
	RAMUsage    string
	CPUUsage    string
	UptimeInfo  string
	IPAddresses []string
	NetworkInfo []string
	Platform    string
	OSInfo      string
}

// GetSystemInfo returns current system information
func GetSystemInfo() SystemInfo {
	info := SystemInfo{}

	// Get platform information
	info.Platform = runtime.GOOS
	hostInfo, err := host.Info()
	if err == nil {
		info.OSInfo = fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.PlatformFamily)
	} else {
		info.OSInfo = "OS information unavailable"
	}

	// Get CPU information
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.CPUInfo = fmt.Sprintf("%d x %s @ %.2f GHz", len(cpuInfo), cpuInfo[0].ModelName, cpuInfo[0].Mhz/1000.0)
	} else {
		info.CPUInfo = "CPU information unavailable"
	}

	// Get CPU usage
	cpuPercent, err := cpu.Percent(100*time.Millisecond, false)
	if err == nil && len(cpuPercent) > 0 {
		info.CPUUsage = fmt.Sprintf("CPU Usage: %2d%%", int(cpuPercent[0]))
	} else {
		info.CPUUsage = "CPU Usage: N/A"
	}

	// Get memory information
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.MemoryInfo = fmt.Sprintf("%d GB System Memory (%.1f GB Used)", memInfo.Total/(1024*1024*1024), float64(memInfo.Used)/(1024*1024*1024))
		info.RAMUsage = fmt.Sprintf("RAM Usage: %2d%%", int(memInfo.UsedPercent))
	} else {
		info.MemoryInfo = "Memory information unavailable"
		info.RAMUsage = "RAM Usage: N/A"
	}

	// Get uptime information
	if err == nil {
		uptime := time.Duration(hostInfo.Uptime) * time.Second
		info.UptimeInfo = fmt.Sprintf("Uptime: %s", formatDuration(uptime))
	} else {
		info.UptimeInfo = "Uptime information unavailable"
	}

	// Get GPU information
	info.GPUInfo = getGPUInfo()

	// Get IP addresses
	info.IPAddresses = getIPAddresses()

	return info
}

// formatDuration formats uptime in a human-readable format
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60

	return fmt.Sprintf("%d hours %02d minutes %02d seconds", hours, mins, secs)
}

// getIPAddresses returns a list of non-loopback IPv4 addresses
func getIPAddresses() []string {
	var addrs []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return []string{"Error getting network interfaces"}
	}

	for _, iface := range ifaces {
		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		ifAddrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range ifAddrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Only include IPv4 addresses
			if ipv4 := ip.To4(); ipv4 != nil {
				addrs = append(addrs, fmt.Sprintf("%s (%s)", ipv4, iface.Name))
			}
		}
	}

	if len(addrs) == 0 {
		return []string{"No IPv4 addresses found"}
	}

	return addrs
}

// getGPUInfo attempts to get GPU information using platform-specific methods
func getGPUInfo() string {
	// Default value
	gpuInfo := "GPU information unavailable"

	// For NVIDIA GPUs on supported platforms
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		cmd := exec.Command("nvidia-smi", "--query-gpu=name", "--format=csv,noheader,nounits")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				gpuInfo = strings.TrimSpace(line)
				break
			}
		}
	} else if runtime.GOOS == "darwin" {
		gpuInfo = "Apple GPU"
	}

	return gpuInfo
}
