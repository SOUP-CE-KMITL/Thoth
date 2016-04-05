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

	Request          int64
	Response         int64
	Response2xx      int64
	Response4xx      int64
	Response5xx      int64
	Response5xxRoute int64
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
	Spec       Spec      `json:"spec"`
	Status     SvcStatus `json:"status"`
}

type Spec struct {
	Ports           []Port   `json:"ports"`
	Selector        Selector `json:"selector"`
	ClusterIP       string   `json:"clusterIP"`
	Type            string   `json:"type"`
	SessionAffinity string   `json:"sessionAffinity"`
}

type SvcSpec struct {
	APIVersion string
	Metadata   struct {
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
	SelfLink          string    `json:"selfLink,omitempty"`
	UID               string    `json:"uid,omitempty"`
	ResourceVersion   string    `json:"resourceVersion,omitempty"`
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`
	Labels            struct {
		Component string `json:"component"`
		Provider  string `json:"provider"`
	} `json:"labels,omitempty"`
}

// Vamp
type Vamp []struct {
	Pxname        string `json:"pxname"`
	Svname        string `json:"svname"`
	Qcur          string `json:"qcur"`
	Qmax          string `json:"qmax"`
	Scur          string `json:"scur"`
	Smax          string `json:"smax"`
	Slim          string `json:"slim"`
	Stot          string `json:"stot"`
	Bin           string `json:"bin"`
	Bout          string `json:"bout"`
	Dreq          string `json:"dreq"`
	Dresp         string `json:"dresp"`
	Ereq          string `json:"ereq"`
	Econ          string `json:"econ"`
	Eresp         string `json:"eresp"`
	Wretr         string `json:"wretr"`
	Wredis        string `json:"wredis"`
	Status        string `json:"status"`
	Weight        string `json:"weight"`
	Act           string `json:"act"`
	Bck           string `json:"bck"`
	Chkfail       string `json:"chkfail"`
	Chkdown       string `json:"chkdown"`
	Lastchg       string `json:"lastchg"`
	Downtime      string `json:"downtime"`
	Qlimit        string `json:"qlimit"`
	Pid           string `json:"pid"`
	Iid           string `json:"iid"`
	Sid           string `json:"sid"`
	Throttle      string `json:"throttle"`
	Lbtot         string `json:"lbtot"`
	Tracked       string `json:"tracked"`
	Rate          string `json:"rate"`
	RateLim       string `json:"rate_lim"`
	RateMax       string `json:"rate_max"`
	CheckStatus   string `json:"check_status"`
	CheckCode     string `json:"check_code"`
	CheckDuration string `json:"check_duration"`
	Hrsp1xx       string `json:"hrsp_1xx"`
	Hrsp2xx       string `json:"hrsp_2xx"`
	Hrsp3xx       string `json:"hrsp_3xx"`
	Hrsp4xx       string `json:"hrsp_4xx"`
	Hrsp5xx       string `json:"hrsp_5xx"`
	HrspOther     string `json:"hrsp_other"`
	Hanafail      string `json:"hanafail"`
	ReqRate       string `json:"req_rate"`
	ReqRateMax    string `json:"req_rate_max"`
	ReqTot        string `json:"req_tot"`
	CliAbrt       string `json:"cli_abrt"`
	SrvAbrt       string `json:"srv_abrt"`
	CompIn        string `json:"comp_in"`
	CompOut       string `json:"comp_out"`
	CompByp       string `json:"comp_byp"`
	CompRsp       string `json:"comp_rsp"`
	Lastsess      string `json:"lastsess"`
	LastChk       string `json:"last_chk"`
	LastAgt       string `json:"last_agt"`
	Qtime         string `json:"qtime"`
	Ctime         string `json:"ctime"`
	Rtime         string `json:"rtime"`
	Ttime         string `json:"ttime"`
}
