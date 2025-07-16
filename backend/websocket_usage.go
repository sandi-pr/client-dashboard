package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/requests"
)

type UsageWebSocket struct {
	upgrader websocket.Upgrader
}

func NewUsageWebSocket() *UsageWebSocket {
	return &UsageWebSocket{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (c *UsageWebSocket) Serve(w http.ResponseWriter, r *http.Request, dataType string) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		interval := 1 * time.Second

		switch dataType {
		case "CPU":
			sendCPUData(conn)
		case "Memory":
			sendMemoryData(conn)
		case "Disk":
			sendDiskData(conn)
		case "GPU":
			sendGPUData(conn)
		default:
			conn.WriteJSON(map[string]string{
				"error": "Unknown data type",
			})
			return
		}

		time.Sleep(interval)
	}
}

func sendCPUData(conn *websocket.Conn) {
	usage, _ := cpu.Percent(0, false)
	info, _ := cpu.Info()
	speed := info[0].Mhz

	procs, _ := process.Processes()
	totalThreads := 0
	for _, p := range procs {
		t, err := p.NumThreads()
		if err == nil {
			totalThreads += int(t)
		}
	}
	uptime, _ := host.Uptime()

	conn.WriteJSON(map[string]interface{}{
		"type":           "cpu_realtime",
		"utilization":    usage[0],
		"speed":          speed,
		"thread":         totalThreads,
		"uptimeReadable": formatUptime(uptime),
		"timestamp":      time.Now().Format(time.RFC3339),
	})
}

func sendMemoryData(conn *websocket.Conn) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting virtual memory:", err)
		return
	}

	swapStat, _ := mem.SwapMemory()

	payload := map[string]interface{}{
		"type":           "memory_realtime",
		"total":          vmStat.Total,
		"used":           vmStat.Used,
		"available":      vmStat.Available,
		"used_percent":   vmStat.UsedPercent,
		"cached":         vmStat.Cached,
		"buffers":        vmStat.Buffers,
		"shared":         vmStat.Shared,
		"active":         vmStat.Active,
		"inactive":       vmStat.Inactive,
		"wired":          vmStat.Wired,
		"committed":      vmStat.Total - vmStat.Available, // approximation
		"paged_pool":     swapStat.Used,                   // not accurate, approximation
		"non_paged_pool": swapStat.Free,                   // not accurate, approximation
		"timestamp":      time.Now().Format(time.RFC3339),
	}

	err = conn.WriteJSON(payload)
	if err != nil {
		log.Println("WebSocket write error:", err)
		return
	}
}

func sendDiskData(conn *websocket.Conn) {
	cmd := exec.Command("iostat", "-dx", "1", "2")
	out, err := cmd.Output()
	if err != nil {
		log.Println("Error executing iostat:", err)
		return
	}

	lines := strings.Split(string(out), "\n")

	startIdx := 0
	headerIdx := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "Device") {
			if headerIdx == -1 {
				headerIdx = i
			} else {
				startIdx = i
				break
			}
		}
	}

	info := map[string]interface{}{
		"type":      "disk_realtime",
		"timestamp": time.Now().Format(time.RFC3339),
		"list":      []map[string]string{},
	}

	for _, line := range lines[startIdx+1:] {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			continue
		}

		device := fields[0]
		if strings.HasPrefix(device, "loop") {
			continue
		}

		d := map[string]string{
			"device":       device,
			"utilization":  fields[13] + " %",
			"read_speed":   fields[5] + " MB/s",
			"write_speed":  fields[6] + " MB/s",
			"avg_response": fields[9] + " ms",
		}

		info["list"] = append(info["list"].([]map[string]string), d)
	}

	err = conn.WriteJSON(info)
	if err != nil {
		log.Println("WebSocket write error:", err)
		return
	}
}

func sendGPUData(conn *websocket.Conn) {
	http := requests.NewHttp(U.Dict{
		"Content-type": "application/json",
	})

	resp, err := http.Post(fmt.Sprintf("%s/api/get-gpu-utilization", USAGE_URL), nil, U.Dict{}.ToJson())

	if err != nil {
		fmt.Println("Error getting GPU utilization:", err)
		return
	}

	if resp.Status != 200 {
		fmt.Println("Error getting GPU utilization:", resp.Status)
		return
	}

	dict := U.NewDictFromJson(resp.Body)

	utilization, _ := dict.GetString("utilization")
	gpuMemoryStr, _ := dict.GetString("gpu_memory")
	dedicatedMemoryStr, _ := dict.GetString("dedicated_gpu_memory")
	temperature, _ := dict.GetString("gpu_temperature")

	var memUsedMB, memTotalMB float64
	fmt.Sscanf(string(gpuMemoryStr), "%f MiB/%f MiB", &memUsedMB, &memTotalMB)

	payload := map[string]interface{}{
		"type":            "gpu_realtime",
		"utilization":     utilization,
		"memory":          fmt.Sprintf("%.1f/%.1f GB", memUsedMB/1024, memTotalMB/1024),
		"dedicatedMemory": dedicatedMemoryStr,
		"temperature":     temperature,
	}

	for i := 0; i < 5; i++ {
		if err := conn.WriteJSON(payload); err != nil {
			fmt.Println("Error sending data to WebSocket:", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func formatUptime(seconds uint64) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	return format2Digit(h) + "h " + format2Digit(m) + "m " + format2Digit(s) + "s"
}

func format2Digit(n uint64) string {
	if n < 10 {
		return "0" + strconv.FormatUint(n, 10)
	}
	return strconv.FormatUint(n, 10)
}
