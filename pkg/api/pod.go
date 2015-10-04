// TODO : need to change to api package
package main

// pod metrics
type Pod struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	Port      int    `json:"port"`
	Ip        string `json:"ip"`
	Cpu       uint64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
	DiskUsage uint64 `json:"disk_usage"`
	Bandwidth uint64 `json:"bandwidth"`
}

// array of pods type
type Pods []Pod
