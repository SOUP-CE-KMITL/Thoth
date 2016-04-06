package profil

import (
	"encoding/json"
	"fmt"
	"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetAllSVC() (thoth.ServiceList, error) {
	response, err := http.Get(thoth.KubeApi + "/api/v1/services")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		panic(err)
	}

	//var objSvc interface{}
	objSvc := thoth.ServiceList{}
	if err := json.Unmarshal([]byte(body), &objSvc); err != nil {
		panic(err)
	}
	return objSvc, err
}

// map[namespace/name]Spec
func GetUserSVC() map[string]thoth.SvcSpec {
	allSvc, err := GetAllSVC()
	if err != nil {
		return nil
	}
	svcMap := make(map[string]thoth.SvcSpec)
	for _, svc := range allSvc.Items {
		if svc.Metadata.Namespace != "default" {
			key := svc.Metadata.Namespace + "/" + svc.Metadata.Name
			svcMap[key] = svc
		}
	}
	return svcMap
}

func GetAllRC() (interface{}, error) {
	response, err := http.Get(thoth.KubeApi + "/api/v1/replicationcontrollers")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		panic(err)
	}

	var objRc interface{}
	if err := json.Unmarshal([]byte(body), &objRc); err != nil {
		panic(err)
	}
	return objRc, err
}

// Get every RC except in default (where kubernetes run)
func GetUserRC() []thoth.RC {
	allRc, err := GetAllRC()
	if err != nil {
		panic(err)
	}
	RCArray := []thoth.RC{}
	_RCLen := len(allRc.(map[string]interface{})["items"].([]interface{}))
	for i := 0; i < _RCLen; i++ {
		namespace := allRc.(map[string]interface{})["items"].([]interface{})[i].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
		if namespace != "default" {
			rc := thoth.RC{
				Name:      allRc.(map[string]interface{})["items"].([]interface{})[i].(map[string]interface{})["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["app"].(string),
				Namespace: namespace,
			}
			fmt.Println(rc.Namespace + "/" + rc.Name)
			RCArray = append(RCArray, rc)
		}
	}
	return RCArray
}

// map[port]{req,rps,2xx,4xx,5xx}
func GetHAProxyStats() map[int]thoth.AppMetric {
	// find the request per sec from haproxy-frontend
	resFront, err := http.Get(thoth.VampApi + "/v1/stats/frontends")
	if err != nil {
		panic(err)
	}
	bodyFront, err := ioutil.ReadAll(resFront.Body)
	resFront.Body.Close()
	if err != nil {
		panic(err)
	}
	//var rps uint64
	var objectFront thoth.Vamp
	err = json.Unmarshal([]byte(bodyFront), &objectFront)
	if err != nil {
		fmt.Println(err)
	}

	//find resonse time from haproxy-backends
	resBack, err := http.Get(thoth.VampApi + "/v1/stats/backends")
	if err != nil {
		panic(err)
	}
	bodyBack, err := ioutil.ReadAll(resBack.Body)
	resBack.Body.Close()
	if err != nil {
		panic(err)
	}
	var objectBack thoth.Vamp
	err = json.Unmarshal([]byte(bodyBack), &objectBack)
	if err != nil {
		fmt.Println(err)
	}
	// Loop Port 9000-9999
	// Get Front
	// Get Back
	// Add to Map

	statsMap := make(map[int]thoth.AppMetric)
	for i := 0; i < len(objectFront); i = i + 2 {
		request, _ := strconv.ParseInt(objectFront[i].ReqRate, 10, 64)
		response, _ := strconv.ParseInt(objectBack[i].Rtime, 10, 64)
		response2xx, _ := strconv.ParseInt(objectBack[i].Hrsp2xx, 10, 64)
		response4xx, _ := strconv.ParseInt(objectBack[i].Hrsp4xx, 10, 64)
		response5xx, _ := strconv.ParseInt(objectBack[i].Hrsp5xx, 10, 64)
		response5xxRoute, _ := strconv.ParseInt(objectBack[i+1].Hrsp5xx, 10, 64)
		met := thoth.AppMetric{
			Request:          request,
			Response:         response,
			Response2xx:      response2xx,
			Response4xx:      response4xx,
			Response5xx:      response5xx,
			Response5xxRoute: response5xxRoute,
		}
		port := 9001 + i/2
		statsMap[port] = met
	}
	// Merge objectFront, objectBack to thoth.AppMetric
	return statsMap
}

//"/api/v1/replicationcontrollers"
//"/api/v1/namespaces/{namespace}/replicationcontrollers"
func GetReplicas(namespace, name string) (int, error) {

	res, err := http.Get(thoth.KubeApi + "/api/v1/namespaces/" + namespace + "/replicationcontrollers/" + name)
	if err != nil {
		//		panic(err)
		return -1, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		//	panic(err)
		return -1, err
	} else {
		// json handler type
		var res_obj map[string]interface{}
		if err := json.Unmarshal(body, &res_obj); err != nil {
			panic(err)
		}
		repli, err := strconv.ParseInt(fmt.Sprint(res_obj["status"].(map[string]interface{})["replicas"]), 10, 32)
		if err != nil {
			panic(err)
		}
		//fmt.Println(json)
		return int(repli), nil
	}
}

/**
	read node resource usage
**/
func GetNodeResource(w http.ResponseWriter, r *http.Request) {
	// get this node memory
	memory, _ := mem.VirtualMemory()
	// get this node cpu percent usage
	cpu_percent, _ := cpu.CPUPercent(time.Duration(1)*time.Second, false)
	// Disk mount Point
	disk_partitions, _ := disk.DiskPartitions(true)
	// Disk usage
	var disk_usages []*disk.DiskUsageStat
	for _, disk_partition := range disk_partitions {
		if disk_partition.Mountpoint == "/" || disk_partition.Mountpoint == "/home" {
			disk_stat, _ := disk.DiskUsage(disk_partition.Device)
			disk_usages = append(disk_usages, disk_stat)
		}
	}
	// Network
	network, _ := net.NetIOCounters(false)

	// create new node obj with resource usage information
	node_metric := thoth.NodeMetric{
		Cpu:       cpu_percent,
		Memory:    memory,
		DiskUsage: disk_usages,
		Network:   network,
	}

	node_json, err := json.MarshalIndent(node_metric, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprint(w, string(node_json))
}

// list all pods
func GetPods() string {
	// to do need to read api and port of api server from configuration file
	res, err := http.Get(thoth.KubeApi + "/api/v1/pods")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	return string(body)
}

/**
 *   Get List of ContainerID and pod's ip by replication name and their namespace
 **/
func GetContainerIDList(url string, rc_name string, namespace string) ([]string, []string, error) {
	// TODO : maybe user want to get container id which map with it's pod
	// initail rasult array
	container_ids := []string{}
	pod_ips := []string{}

	res, err := http.Get(url + "/api/v1/namespaces/" + namespace + "/pods")
	if err != nil {
		fmt.Println("Can't connect to cadvisor")
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, nil, err
	} else {
		// json handler type
		var res_obj map[string]interface{}
		if err := json.Unmarshal(body, &res_obj); err != nil {
			return nil, nil, err
		}
		pod_arr := res_obj["items"].([]interface{})
		// iterate to get pod of specific rc
		for _, pod := range pod_arr {
			pod_name := pod.(map[string]interface{})["metadata"].(map[string]interface{})["generateName"]
			if pod_name != nil {
				if pod_name == rc_name+"-" {
					pod_ip := pod.(map[string]interface{})["status"].(map[string]interface{})["podIP"]
					if pod_ip != nil {
						pod_ips = append(pod_ips, pod_ip.(string))
						containers := pod.(map[string]interface{})["status"].(map[string]interface{})["containerStatuses"].([]interface{})
						// one pod can has many container ,so iterate for get each container
						for _, container := range containers {
							container_id := container.(map[string]interface{})["containerID"].(string)[9:]
							container_ids = append(container_ids, container_id)
						}
					}
				}
			}
		}
		return container_ids, pod_ips, nil
	}
}

/**
 CPU Percent Calculation
**/
func DockerCPUPercent(container_id string) (float32, error) {
	res, err := GetCpu(container_id)
	if err != nil {
		return 0.0, nil
	}
	return float32(res), nil

}

/**
 	get resource usage of application (pods) on node
**/
func GetAppResource(namespace, name string) thoth.AppMetric {
	var summary_cpu float32
	var memory_bundle []*docker.CgroupMemStat

	container_ids, pod_ips, err := GetContainerIDList(thoth.KubeApi, name, namespace)
	if err != nil {
		panic(err)
	}
	for _, container_id := range container_ids {
		//	fmt.Println(container_id, pod_ips)
		// calculation percentage of cpu usage
		container_cpu, _ := DockerCPUPercent(container_id)
		summary_cpu += container_cpu
		// memory usage
		container_memory, _ := docker.CgroupMemDocker(container_id)
		memory_bundle = append(memory_bundle, container_memory)
	}

	podNum := len(pod_ips)
	/*
	 */
	//fmt.Println("rps: ", rps, ", rtime: ", rtime)
	// find the cpu avarage of application cpu usage
	average_cpu := summary_cpu / float32(len(container_ids))
	// Cal Avg Mem usage
	var avgMem uint64
	for i := 0; i < podNum; i++ {
		if memory_bundle[i] != nil {
			avgMem += memory_bundle[i].MemUsageInBytes
		}
	}
	avgMem = avgMem / uint64(podNum)
	avgMem = avgMem / uint64(1024*1024) // MB
	// Convert to percentage 100%=200MB
	avgMem = avgMem / 2

	// create appliction object
	app_metric := thoth.AppMetric{
		App:    name,
		Cpu:    average_cpu,
		Memory: int64(avgMem),
		/*
			Request:     rps_int,
			Response:    rtime_int,
			Response2xx: res2xx_int,
			Response4xx: res4xx_int,
			Response5xx: res5xx_int,
		*/
	}
	return app_metric

}

func GetCpu(containerId string) (float64, error) {
	var err error
	var result []byte
	if result, err = exec.Command("docker", "stats", "--no-stream", containerId).Output(); err != nil {
		panic(err)
	}

	cpuPercent := strings.Fields(string(result))[10]
	//cpuPercent := strings.Fields(string(result))[14]
	//fmt.Println("get CPU ", cpuPercent2)
	cpuPercent = cpuPercent[0 : len(cpuPercent)-1]
	return strconv.ParseFloat(cpuPercent, 32)
}
