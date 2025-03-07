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
	return SystemInfo{
		CPUInfo:     GetCPUInfo(),
		GPUInfo:     GetGPUInfo(),
		MemoryInfo:  GetMemoryInfo(),
		RAMUsage:    GetRAMUsage(),
		CPUUsage:    GetCPUUsage(),
		UptimeInfo:  GetUptimeInfo(),
		IPAddresses: GetIPAddresses(),
		CPUTemp:     GetCPUTemperature(),
		GPUTemp:     GetGPUTemperature(),
	}
}

// GetCPUInfo returns CPU information
func GetCPUInfo() string {
	cpuInfo, err := cpu.Info()
	if err != nil || len(cpuInfo) == 0 {
		return "CPU information unavailable"
	}
	return fmt.Sprintf("%d x %s @ %.2f GHz", len(cpuInfo), cpuInfo[0].ModelName, cpuInfo[0].Mhz/1000.0)
}

// GetGPUInfo returns GPU information
func GetGPUInfo() string {
	// In a real implementation, you would use platform-specific methods to get GPU info
	// For example, on Linux you might parse the output of `lspci` or use a library
	// For demonstration purposes, we'll return a placeholder
	return "NVIDIA GeForce RTX 3080 (10GB VRAM)"
}

// GetCPUUsage returns current CPU usage percentage
func GetCPUUsage() string {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil || len(cpuPercent) == 0 {
		return "CPU Usage: N/A"
	}
	return fmt.Sprintf("CPU Usage: %.1f%%", cpuPercent[0])
}

// GetMemoryInfo returns memory information
func GetMemoryInfo() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "Memory information unavailable"
	}
	return fmt.Sprintf("%.1f GB System Memory", float64(memInfo.Total)/(1024*1024*1024))
}

// GetRAMUsage returns RAM usage percentage
func GetRAMUsage() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "RAM Usage: N/A"
	}
	return fmt.Sprintf("RAM Usage: %.1f%%", memInfo.UsedPercent)
}

// GetUptimeInfo returns system uptime
func GetUptimeInfo() string {
	hostInfo, err := host.Info()
	if err != nil {
		return "Uptime information unavailable"
	}
	uptime := time.Duration(hostInfo.Uptime) * time.Second
	return fmt.Sprintf("System uptime: %s", FormatDuration(uptime))
}

// FormatDuration formats uptime in a human-readable format
func FormatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)
	}
	return fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
}

// GetIPAddresses returns a list of non-loopback IPv4 addresses
func GetIPAddresses() []string {
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

// GetCPUTemperature returns CPU temperature
func GetCPUTemperature() string {
	// On macOS, we can use the "smc" command line tool to get CPU temperature
	// On Linux, we can read from /sys/class/thermal/thermal_zone*/temp
	// For now, we'll return a placeholder

	// In a real implementation, you would use platform-specific methods to get the actual temperature
	// For example, on Linux:
	// temperatures, err := host.SensorsTemperatures()
	// if err == nil {
	//     for _, temp := range temperatures {
	//         if strings.Contains(strings.ToLower(temp.SensorKey), "cpu") {
	//             return fmt.Sprintf("CPU: %.1f°C", temp.Temperature)
	//         }
	//     }
	// }

	// For demonstration purposes, we'll simulate a temperature
	return fmt.Sprintf("CPU Temp: %.1f°C", 45.0+5.0*float64(time.Now().Second()%10)/10.0)
}

// GetGPUTemperature returns GPU temperature
func GetGPUTemperature() string {
	// Similar to CPU temperature, this would be platform-specific
	// For demonstration purposes, we'll simulate a temperature
	return fmt.Sprintf("GPU Temp: %.1f°C", 60.0+10.0*float64(time.Now().Second()%10)/10.0)
}
