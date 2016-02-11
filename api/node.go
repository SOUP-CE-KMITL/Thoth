// TODO : need to change to api package
package main
import (
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/net"
)

// node metrics
type Node struct {
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	Cpu       float64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
	DiskUsage uint64 `json:"disk_usage"`
}

type NodeMetric struct {
	Cpu       []float64 `json:"cpu"`
	Memory    *mem.VirtualMemoryStat `json:"memory"`
	DiskUsage []*disk.DiskUsageStat `json:"disk_usage"`
	Network   []net.NetIOCountersStat `json:"network"`
}

// array of node type
type Nodes []Node
