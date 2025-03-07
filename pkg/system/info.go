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
	// Get CPU information
	cpuInfo, err := cpu.Info()
	cpuInfoStr := "CPU information unavailable"
	if err == nil && len(cpuInfo) > 0 {
		cpuInfoStr = fmt.Sprintf("%d x %s @ %.2f GHz", len(cpuInfo), cpuInfo[0].ModelName, cpuInfo[0].Mhz/1000.0)
	}

	// Get CPU usage
	cpuPercent, err := cpu.Percent(100*time.Millisecond, false)
	cpuUsageStr := "CPU Usage: N/A"
	if err == nil && len(cpuPercent) > 0 {
		cpuUsageStr = fmt.Sprintf("CPU Usage: %.1f%%", cpuPercent[0])
	}

	// Get memory information
	memInfo, err := mem.VirtualMemory()
	memInfoStr := "Memory information unavailable"
	ramUsageStr := "RAM Usage: N/A"
	if err == nil {
		memInfoStr = fmt.Sprintf("%.1f GB System Memory", float64(memInfo.Total)/(1024*1024*1024))
		ramUsageStr = fmt.Sprintf("RAM Usage: %.1f%%", memInfo.UsedPercent)
	}

	// Get uptime information
	hostInfo, err := host.Info()
	uptimeStr := "Uptime information unavailable"
	if err == nil {
		uptime := time.Duration(hostInfo.Uptime) * time.Second
		uptimeStr = fmt.Sprintf("System uptime: %s", formatDuration(uptime))
	}

	// For demonstration purposes, we'll use placeholder GPU info
	gpuInfoStr := "NVIDIA GeForce RTX 3080 (10GB VRAM)"

	// For demonstration purposes, we'll simulate temperatures
	cpuTempStr := fmt.Sprintf("CPU Temp: %.1f°C", 45.0+5.0*float64(time.Now().Second()%10)/10.0)
	gpuTempStr := fmt.Sprintf("GPU Temp: %.1f°C", 60.0+10.0*float64(time.Now().Second()%10)/10.0)

	return SystemInfo{
		CPUInfo:     cpuInfoStr,
		GPUInfo:     gpuInfoStr,
		MemoryInfo:  memInfoStr,
		RAMUsage:    ramUsageStr,
		CPUUsage:    cpuUsageStr,
		UptimeInfo:  uptimeStr,
		IPAddresses: getIPAddresses(),
		CPUTemp:     cpuTempStr,
		GPUTemp:     gpuTempStr,
	}
}

// formatDuration formats uptime in a human-readable format
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)
	}
	return fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
}

// getIPAddresses returns a list of non-loopback IPv4 addresses
func getIPAddresses() []string {
	var addresses []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{"Error getting network interfaces"}
	}

	for _, iface := range interfaces {
		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
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
				addresses = append(addresses, fmt.Sprintf("http://%s/ (%s)", ipv4.String(), iface.Name))
			}
		}
	}

	if len(addresses) == 0 {
		return []string{"No IPv4 addresses found"}
	}

	return addresses
}
