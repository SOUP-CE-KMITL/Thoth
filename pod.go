// TODO : need to change to api package
package thoth

import (
	"time"
)

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

type PodList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []PodItem `json:"items"`
}

type PodItem struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Labels            struct {
			App string `json:"app"`
		} `json:"labels"`
		Annotations struct {
			KubernetesIoConfigHash   string    `json:"kubernetes.io/config.hash"`
			KubernetesIoConfigMirror string    `json:"kubernetes.io/config.mirror"`
			KubernetesIoConfigSeen   time.Time `json:"kubernetes.io/config.seen"`
			KubernetesIoConfigSource string    `json:"kubernetes.io/config.source"`
		} `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Containers []struct {
			Name      string   `json:"name"`
			Image     string   `json:"image"`
			Command   []string `json:"command"`
			Resources struct {
			} `json:"resources"`
			TerminationMessagePath string `json:"terminationMessagePath"`
			ImagePullPolicy        string `json:"imagePullPolicy"`
		} `json:"containers"`
		RestartPolicy                 string `json:"restartPolicy"`
		TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
		DNSPolicy                     string `json:"dnsPolicy"`
		NodeName                      string `json:"nodeName"`
		HostNetwork                   bool   `json:"hostNetwork"`
	} `json:"spec"`
	Status struct {
		Phase      string `json:"phase"`
		Conditions []struct {
			Type               string      `json:"type"`
			Status             string      `json:"status"`
			LastProbeTime      interface{} `json:"lastProbeTime"`
			LastTransitionTime interface{} `json:"lastTransitionTime"`
		} `json:"conditions"`
		HostIP            string    `json:"hostIP"`
		PodIP             string    `json:"podIP"`
		StartTime         time.Time `json:"startTime"`
		ContainerStatuses []struct {
			Name  string `json:"name"`
			State struct {
				Running struct {
					StartedAt time.Time `json:"startedAt"`
				} `json:"running"`
			} `json:"state"`
			LastState struct {
			} `json:"lastState"`
			Ready        bool   `json:"ready"`
			RestartCount int    `json:"restartCount"`
			Image        string `json:"image"`
			ImageID      string `json:"imageID"`
			ContainerID  string `json:"containerID"`
		} `json:"containerStatuses"`
	} `json:"status"`
}
