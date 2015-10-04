// TODO : need to change to api package
package main

// node metrics
type Node struct {
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	Cpu       uint64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
	DiskUsage uint64 `json:"disk_usage"`
}

// array of node type
type Nodes []Node
