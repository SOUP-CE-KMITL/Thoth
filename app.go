// TODO : need to change to api package
package thoth

import (
	"time"
)

var KubeApi string = "http://localhost:8080"
var InfluxdbApi string = "http://localhost:8086"
var VampApi string = "http://localhost:10001"

type AppMetric struct {
	App    string  `json:"app"`
	Cpu    float32 `json:"cpu"`
	Memory int64   `json:"memory"`
	/*
		Request     int64
		Response    int64
		Response2xx int64
		Response4xx int64
		Response5xx int64
	*/
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

type SvcSpec struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Spec struct {
		Ports []struct {
			Protocol   string `json:"protocol"`
			Port       int    `json:"port"`
			TargetPort int    `json:"targetPort"`
			NodePort   int    `json:"nodePort"`
		} `json:"ports"`
		Selector struct {
			App string `json:"app"`
		} `json:"selector"`
		ClusterIP       string `json:"clusterIP"`
		Type            string `json:"type"`
		SessionAffinity string `json:"sessionAffinity"`
	} `json:"spec"`
	Status struct {
		LoadBalancer struct {
		} `json:"loadBalancer"`
	} `json:"status"`
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

type ServiceList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []SvcSpec `json:"items"`
}
type Metadata struct {
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	SelfLink          string    `json:"selfLink"`
	UID               string    `json:"uid"`
	ResourceVersion   string    `json:"resourceVersion"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
	Labels            struct {
		Component string `json:"component"`
		Provider  string `json:"provider"`
	} `json:"labels"`
}

// Vamp Frontend
type VampFrontend []struct {
	Name         string `json:"name"`
	Mode         string `json:"mode"`
	BindPort     int    `json:"bindPort"`
	BindIP       string `json:"bindIp"`
	UnixSock     string `json:"unixSock"`
	SockProtocol string `json:"sockProtocol"`
	Options      struct {
		AbortOnClose    bool `json:"abortOnClose"`
		AllBackups      bool `json:"allBackups"`
		CheckCache      bool `json:"checkCache"`
		ForwardFor      bool `json:"forwardFor"`
		HTTPClose       bool `json:"httpClose"`
		HTTPCheck       bool `json:"httpCheck"`
		SslHelloCheck   bool `json:"sslHelloCheck"`
		TCPKeepAlive    bool `json:"tcpKeepAlive"`
		TCPLog          bool `json:"tcpLog"`
		TCPSmartAccept  bool `json:"tcpSmartAccept"`
		TCPSmartConnect bool `json:"tcpSmartConnect"`
	} `json:"options"`
	DefaultBackend string `json:"defaultBackend"`
	HTTPQuota      struct {
	} `json:"httpQuota"`
	TCPQuota struct {
	} `json:"tcpQuota"`
}

// Vamp Backend
type VampBackend []struct {
	Name    string `json:"name"`
	Mode    string `json:"mode"`
	Servers []struct {
		Name          string `json:"name"`
		Host          string `json:"host"`
		Port          int    `json:"port"`
		UnixSock      string `json:"unixSock"`
		Weight        int    `json:"weight"`
		Maxconn       int    `json:"maxconn"`
		Check         bool   `json:"check"`
		CheckInterval int    `json:"checkInterval"`
	} `json:"servers"`
	Options struct {
		AbortOnClose    bool `json:"abortOnClose"`
		AllBackups      bool `json:"allBackups"`
		CheckCache      bool `json:"checkCache"`
		ForwardFor      bool `json:"forwardFor"`
		HTTPClose       bool `json:"httpClose"`
		HTTPCheck       bool `json:"httpCheck"`
		SslHelloCheck   bool `json:"sslHelloCheck"`
		TCPKeepAlive    bool `json:"tcpKeepAlive"`
		TCPLog          bool `json:"tcpLog"`
		TCPSmartAccept  bool `json:"tcpSmartAccept"`
		TCPSmartConnect bool `json:"tcpSmartConnect"`
	} `json:"options"`
	ProxyMode bool `json:"proxyMode"`
}
