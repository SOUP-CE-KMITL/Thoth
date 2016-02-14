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
	"strconv"
	"time"
)

//"/api/v1/replicationcontrollers"
//"/api/v1/namespaces/{namespace}/replicationcontrollers"
func GetReplicas(namespace, name string) int {

	res, err := http.Get(thoth.KubeApi + "/api/v1/namespaces/" + namespace + "/replicationcontrollers/" + name)
	if err != nil {
		fmt.Println("Can't connect to cadvisor")
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
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
		return int(repli)
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
					pod_ips = append(pod_ips, pod.(map[string]interface{})["status"].(map[string]interface{})["podIP"].(string))
					containers := pod.(map[string]interface{})["status"].(map[string]interface{})["containerStatuses"].([]interface{})
					// one pod can has many container ,so iterate for get each container
					for _, container := range containers {
						container_id := container.(map[string]interface{})["containerID"].(string)[9:]
						container_ids = append(container_ids, container_id)
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
func DockerCPUPercent(interval time.Duration, container_id string) (float64, error) {
	getAllBusy := func(t *cpu.CPUTimesStat) (float64, float64) {
		busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
			t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Stolen
		return busy + t.Idle, busy
	}

	calculate := func(t1, t2 *cpu.CPUTimesStat) float64 {
		t1All, t1Busy := getAllBusy(t1)
		t2All, t2Busy := getAllBusy(t2)
		if t2Busy <= t1Busy {
			return 0
		}
		if t2All <= t1All {
			return 1
		}
		//	fmt.Println("Busy: ", t2Busy-t1Busy, ", All: ", t2All-t1All)
		//	fmt.Println("idle time 1: ", t1Busy-t1All, ", idle time 2: ", t2Busy-t2All)
		return (t2Busy - t1Busy) / (t2All - t1All) * 100
	}

	// Get CPU usage at the start of the interval.
	var cpuTimes1 *cpu.CPUTimesStat
	cpuTimes1, _ = docker.CgroupCPUDocker(container_id)

	if interval > 0 {
		time.Sleep(interval)
	}

	// And at the end of the interval.
	var cpuTimes2 *cpu.CPUTimesStat
	cpuTimes2, _ = docker.CgroupCPUDocker(container_id)

	var rets float64
	rets = calculate(cpuTimes1, cpuTimes2)
	return rets, nil
}

/**
 	get resource usage of application (pods) on node
**/
func GetAppResource(namespace, name string) thoth.AppMetric {
	var summary_cpu float64
	var memory_bundle []*docker.CgroupMemStat

	container_ids, pod_ips, err := GetContainerIDList(thoth.KubeApi, name, namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, container_id := range container_ids {
		//	fmt.Println(container_id, pod_ips)
		// calculation percentage of cpu usage
		container_cpu, _ := DockerCPUPercent(time.Duration(1)*time.Second, container_id)
		summary_cpu += container_cpu
		// memory usage
		container_memory, _ := docker.CgroupMemDocker(container_id)
		memory_bundle = append(memory_bundle, container_memory)
	}

	podNum := len(pod_ips)

	// find the request per sec from haproxy-frontend
	res_front, err := http.Get(thoth.VampApi + "/v1/stats/frontends")
	if err != nil {
		panic(err)
	}
	body_front, err := ioutil.ReadAll(res_front.Body)
	res_front.Body.Close()
	if err != nil {
		panic(err)
	}
	//var rps uint64
	var object_front []map[string]interface{}
	err = json.Unmarshal([]byte(body_front), &object_front)
	rps := object_front[0]["req_rate"].(string)
	rps_int, _ := strconv.ParseInt(rps, 10, 64)
	if err == nil {
	} else {
		fmt.Println(err)
	}

	//find resonse time from haproxy-backends
	//var rtime uint64
	res_back, err := http.Get(thoth.VampApi + "/v1/stats/backends")
	if err != nil {
		panic(err)
	}
	body_back, err := ioutil.ReadAll(res_back.Body)
	res_back.Body.Close()
	if err != nil {
		panic(err)
	}

	var object_back []map[string]interface{}
	err = json.Unmarshal([]byte(body_back), &object_back)
	rtime := object_back[0]["rtime"].(string)
	res_2xx := object_back[0]["hrsp_2xx"].(string)
	res_4xx := object_back[0]["hrsp_4xx"].(string)
	res_5xx := object_back[0]["hrsp_5xx"].(string)
	rtime_int, _ := strconv.ParseInt(rtime, 10, 64)
	res2xx_int, _ := strconv.ParseInt(res_2xx, 10, 64)
	res4xx_int, _ := strconv.ParseInt(res_4xx, 10, 64)
	res5xx_int, _ := strconv.ParseInt(res_5xx, 10, 64)
	if err == nil {
	} else {
		fmt.Println(err)
	}

	//fmt.Println("rps: ", rps, ", rtime: ", rtime)
	// find the cpu avarage of application cpu usage
	average_cpu := summary_cpu / float64(len(container_ids))
	// Cal Avg Mem usage
	var avgMem uint64
	for i := 0; i < podNum; i++ {
		avgMem += memory_bundle[i].MemUsageInBytes
	}
	avgMem = avgMem / uint64(podNum)
	avgMem = avgMem / uint64(1024*1024) // MB

	// create appliction object
	app_metric := thoth.AppMetric{
		App:         name,
		Cpu:         average_cpu,
		Memory:      int64(avgMem),
		Request:     rps_int,
		Response:    rtime_int,
		Response2xx: res2xx_int,
		Response4xx: res4xx_int,
		Response5xx: res5xx_int,
	}
	return app_metric

}
