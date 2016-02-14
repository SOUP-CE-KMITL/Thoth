// TODO : need to change to api package
package thoth

type AppMetric struct {
	App         string  `json:"app"`
	Cpu         float32 `json:"cpu"`
	Memory      int64   `json:"memory"`
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
	Name      string `json: "name"`
	Replicas  int    `json: "replicas"`
	Namespace string `json: "namespace"`
	Image     string `json: "image"`
	Port      int    `json: "port"`
}
type RC struct {
	Namespace string // Namespace = User
	Name      string
}

type Svc struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Metadata   Metadata  `json:"metadata"`
	Spec       SvcSpec   `json:"spec"`
	Status     SvcStatus `json:"status"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
type SvcSpec struct {
	Ports    []Port   `json:"ports"`
	Selector Selector `json:"selector"`
	Type     string   `json:"type"`
}
type Port struct {
	NodePort int `json:"nodePort"`
	Port     int `json:"port"`
	//	Protocol   string `json:"protocol"`
	TargetPort int `json:"targetPort"`
}
type Selector struct {
	App string `json:"app"`
}
type SvcStatus struct {
	LoadBalancer SvcLoadBalancer `json:"loadBalancer"`
}

type SvcLoadBalancer struct {
	Ingress []SvcIngress `json:"ingress"`
}
type SvcIngress struct {
	IP string `json:"ip"`
}

var KubeApi string = "http://localhost:8080"
var InfluxdbApi string = "http://localhost:8086"
var VampApi string = "http://localhost:10001"
