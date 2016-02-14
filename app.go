// TODO : need to change to api package
package thoth

import (
	"github.com/shirou/gopsutil/docker"
)

type AppMetric struct {
	App         string                  `json:"app"`
	Cpu         float64                 `json:"cpu"`
	Memory      []*docker.CgroupMemStat `json:"memory"`
	Request     int64
	Response    int64
	Response2xx int64
	Response4xx int64
	Response5xx int64
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


// Replication 

type SendRC struct {
	Name string `json: "name"`
	Replicas int `json: "replicas"`
	Namespace string `json: "namespace"`
	Image string `json: "image"`
	Port int `json: "port"`
}