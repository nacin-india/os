package system

import (
	"fmt"
	"net"
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
	CPUTemp     string
	GPUTemp     string
}

// GetSystemInfo returns current system information
func GetSystemInfo() SystemInfo {
	info := SystemInfo{}

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
		info.CPUUsage = fmt.Sprintf("CPU Usage: %02d%%", int(cpuPercent[0]))
	} else {
		info.CPUUsage = "CPU Usage: N/A"
	}

	// Get memory information
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.MemoryInfo = fmt.Sprintf("%.1f GB System Memory", float64(memInfo.Total)/(1024*1024*1024))
		info.RAMUsage = fmt.Sprintf("RAM Usage: %02d%%", int(memInfo.UsedPercent))
	} else {
		info.MemoryInfo = "Memory information unavailable"
		info.RAMUsage = "RAM Usage: N/A"
	}

	// Get uptime information
	hostInfo, err := host.Info()
	if err == nil {
		uptime := time.Duration(hostInfo.Uptime) * time.Second
		info.UptimeInfo = fmt.Sprintf("System uptime: %s", formatDuration(uptime))
	} else {
		info.UptimeInfo = "Uptime information unavailable"
	}

	// Simulated values for demo
	info.GPUInfo = "NVIDIA GeForce RTX 3080 (10GB VRAM)"
	info.CPUTemp = fmt.Sprintf("CPU Temp: %d°C", int(45.0+5.0*float64(time.Now().Second()%10)/10.0))
	info.GPUTemp = fmt.Sprintf("GPU Temp: %d°C", int(60.0+10.0*float64(time.Now().Second()%10)/10.0))

	// Get IP addresses
	info.IPAddresses = getIPAddresses()

	return info
}

// formatDuration formats uptime in a human-readable format
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %02d hours %02d minutes %02d seconds", days, hours, mins, secs)
	}
	return fmt.Sprintf("%02d hours %02d minutes %02d seconds", hours, mins, secs)
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
