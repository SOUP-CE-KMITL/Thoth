// TODO : need to change to api package
package main

type AppMetrics struct {
	App       uint64 `json:"app"`
	Cpu       uint64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
	DiskUsage uint64 `json:"disk_usage"`
	Bandwidth uint64 `json:"bandwidth"`
}

// This is schema of Application Profile
type App struct {
	Name       string `json:"name"`
	ExternalIp string `json:"external_ip"`
	InternalIp string `json:"internal_ip"`
	Image      string `json:"image"`
	Pods       []Pod  `json:"pods"`
}

type Apps []App
