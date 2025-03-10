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
	info := SystemInfo{
		Platform:   runtime.GOOS,
		OSInfo:     "OS information unavailable",
		CPUInfo:    "CPU information unavailable",
		CPUUsage:   "CPU Usage: N/A",
		MemoryInfo: "Memory information unavailable",
		RAMUsage:   "RAM Usage: N/A",
		UptimeInfo: "Uptime information unavailable",
	}

	// Get platform and host information
	if hostInfo, err := host.Info(); err == nil {
		info.OSInfo = fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.PlatformFamily)
		uptime := time.Duration(hostInfo.Uptime) * time.Second
		hours, mins, secs := int(uptime.Hours()), int(uptime.Minutes())%60, int(uptime.Seconds())%60
		info.UptimeInfo = fmt.Sprintf("Uptime: %d hours %02d minutes %02d seconds", hours, mins, secs)
	}

	// Get CPU information
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		info.CPUInfo = fmt.Sprintf("%d x %s @ %.2f GHz", len(cpuInfo), cpuInfo[0].ModelName, cpuInfo[0].Mhz/1000.0)
	}

	// Get CPU usage
	if cpuPercent, err := cpu.Percent(100*time.Millisecond, false); err == nil && len(cpuPercent) > 0 {
		info.CPUUsage = fmt.Sprintf("CPU Usage: %2d%%", int(cpuPercent[0]))
	}

	// Get memory information
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info.MemoryInfo = fmt.Sprintf("%.1f GB System Memory (%.2f GB Used)",
			float64(memInfo.Total)/float64(1000*1000*1000),
			float64(memInfo.Used)/float64(1000*1000*1000),
		)
		info.RAMUsage = fmt.Sprintf("RAM Usage: %2d%%", int(memInfo.UsedPercent))
	}

	// Get GPU and IP information
	info.GPUInfo = getGPUInfo()
	info.IPAddresses = getIPAddresses()

	return info
}

// getIPAddresses returns a list of non-loopback IPv4 addresses
func getIPAddresses() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []string{"Error getting network interfaces"}
	}

	var addrs []string
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

			// Only include non-loopback IPv4 addresses
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				addrs = append(addrs, fmt.Sprintf("%s (%s)", ip.To4(), iface.Name))
			}
		}
	}

	if len(addrs) == 0 {
		return []string{"No IP addresses found"}
	}
	return addrs
}

// getGPUInfo attempts to get GPU information using platform-specific methods
func getGPUInfo() string {
	switch runtime.GOOS {
	case "darwin":
		return "Apple GPU"
	case "linux", "windows":
		cmd := exec.Command("nvidia-smi", "--query-gpu=name", "--format=csv,noheader,nounits")
		if output, err := cmd.Output(); err == nil {
			for _, line := range strings.Split(string(output), "\n") {
				if line = strings.TrimSpace(line); line != "" {
					return line
				}
			}
		}
	}
	return "GPU information unavailable"
}
