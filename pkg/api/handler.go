// TODO : need to change to api
package main

import (
	"bytes"
	"encoding/json"
	"encoding/binary"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os/exec"
	"errors"
	"strconv"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"time"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/docker"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// list every node
func GetNodes(w http.ResponseWriter, r *http.Request) {
	// to do need to read api and port of api server from configuration file
	res, err := http.Get("http://localhost:8080/api/v1/nodes")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	// defer for ensure that res is close.
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))
}

func GetNode(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// node name from user.
	nodesName := vars["nodeName"]
	// TODO: need to read api and port of api server from configuration file
	res, err := http.Get("http://localhost:8080/api/v1/nodes/" + nodesName)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	var object map[string]interface{}
	err = json.Unmarshal([]byte(body), &object)
	if err == nil {
		fmt.Printf("%+v\n", object)
	} else {
		fmt.Println(err)
	}
	send_obj, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprint(w, string(send_obj))
}

func OptionCors(w http.ResponseWriter, r *http.Request) {
	// TODO: need to change origin to deployed domain name
	if origin := r.Header.Get("Origin"); origin != "http://localhost" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
}

// list specific node cpu
func NodeCpu(w http.ResponseWriter, r *http.Request) {
}

// list specifc node memory
func NodeMemory(w http.ResponseWriter, r *http.Request) {
}

// list all pods
func GetPods(w http.ResponseWriter, r *http.Request) {
	// to do need to read api and port of api server from configuration file
	res, err := http.Get("http://localhost:8080/api/v1/pods")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(body))
}

// list specific pod details
func GetPod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// node name from user.
	podName := vars["podName"]
	fmt.Fprint(w, string(podName))
	// to do need to read api and port of api server from configuration file
	// TODO: change namespace to flexible.
	res, err := http.Get("http://localhost:8080/api/v1/namespaces/default/pods/" + podName)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(body))
}

// list specific pod cpu
func PodCpu(w http.ResponseWriter, r *http.Request) {
}

// list specific pod memory
func PodMemory(w http.ResponseWriter, r *http.Request) {
}

// post handler for scale pod by pod name

// TODO : remove
// test mocks
func nodeTestMock(w http.ResponseWriter, r *http.Request) {
	nodes := Nodes{
		Node{Name: "node1", Ip: "192.168.1.2", Cpu: 5000, Memory: 3000, DiskUsage: 1000},
		Node{Name: "node2", Ip: "192.169.1.4", Cpu: 5000, Memory: 3000, DiskUsage: 1000},
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		panic(err)
	}
}

// TODO : remove
// test ssh to exec command on other machine
func testExec(w http.ResponseWriter, r *http.Request) {
	commander := SSHCommander{"root", "161.246.70.75"}
	cmd := []string{
		"ls",
		".",
	}
	var (
		output []byte
		err    error
	)

	if output, err = commander.Command(cmd...).Output(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(output[:6]))
}

func GetApp(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	// app name from user.
	appName := vars["appName"]
	fmt.Println(appName)
	// TODO: need to find new solution to get info from api like other done.
	res, err := exec.Command("kubectl", "get", "pod", "-l", "app="+appName, "-o", "json").Output()
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(res))
}


func CreatePod(w http.ResponseWriter, r *http.Request) {
	var pod Pod
	// limits json post request for prevent overflow attack.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	// catch error from close reader
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// get request json information
	if err := json.Unmarshal(body, &pod); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// prepare json to send to create by kubernetes api server
	labels := map[string]interface{}{
		"app": pod.Name,
	}

	metadata := map[string]interface{}{
		"name":   pod.Name,
		"labels": labels,
	}

	ports := map[string]interface{}{
		"containerPort": 80,
	}

	containers := map[string]interface{}{
		"name":   pod.Name,
		"image":  pod.Image,
		"ports":  []map[string]interface{}{ports},
		"memory": pod.Memory,
		"cpu":    pod.Cpu,
	}

	spec := map[string]interface{}{
		"containers": []map[string]interface{}{containers},
	}

	objReq := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata":   metadata,
		"spec":       spec,
	}

	jsonReq, err := json.MarshalIndent(objReq, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println("you sent ", string(jsonReq))
	// post json to kubernete api server

	// TODO: need to change name space to user namespace
	postUrl := "http://localhost:8080/api/v1/namespaces/default/pods"
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonReq))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// defer for ensure
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(response))
}	

/**
* 	Get current resource by using cgroup
**/
func GetCurrentResourceCgroup(container_id string, metric_type int) (uint64, error) {
	// TODO : Read Latency from HA Proxy
	// file path prefix
	var path = "/sys/fs/cgroup/memory/docker/" + container_id + "/"

	if metric_type == 2 {
		// read memory usage
		current_usage, err := ioutil.ReadFile(path + "memory.usage_in_bytes")
		if err != nil {
			return binary.BigEndian.Uint64(current_usage), nil
		} else {
			return 0, err
		}
	} else if metric_type == 3 {
		// read memory usage
		current_usage, err := ioutil.ReadFile(path + "memory.usage_in_bytes")
		if err != nil {
			return 0, err
		} else {
			n := bytes.Index(current_usage, []byte{10})
			usage_str := string(current_usage[:n])
			resource_usage, _ := strconv.ParseInt(usage_str, 10, 64)
			return uint64(resource_usage), nil
		}
	} else {
		// not match any case
		return 0, errors.New("not match any case")
	}
}

/**
 *   Get List of ContainerID and pod's ip by replication name and their namespace
 **/
func GetContainerIDList(url string, port string, rc_name string, namespace string) ([]string, []string, error) {
	// TODO : maybe user want to get container id which map with it's pod
	// initail rasult array
	container_ids := []string{}
	pod_ips := []string{}

	res, err := http.Get(url + ":" + port + "/api/v1/namespaces/" + namespace + "/pods")
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
	read node resource usage
**/
func GetNodeResource(w http.ResponseWriter, r *http.Request) {
	// get this node memory 
	memory, _ := mem.VirtualMemory()
	// get this node cpu percent usage
	cpu_percent, _ := cpu.CPUPercent(time.Duration(1) * time.Second, false)
	// Disk mount Point
	disk_partitions, _ := disk.DiskPartitions(true)
	// Disk usage
	var disk_usages []*disk.DiskUsageStat
	for _, disk_partition := range disk_partitions {
		if disk_partition.Mountpoint == "/" || disk_partition.Mountpoint == "/home" {
			disk_stat, _ := disk.DiskUsage(disk_partition.Device);
			disk_usages = append(disk_usages, disk_stat)
		} 
	}
	// Network
	network, _ := net.NetIOCounters(false)

	// create new node obj with resource usage information
	node_metric := NodeMetric{
		Cpu: cpu_percent,
		Memory: memory,
		DiskUsage: disk_usages,
		Network: network,
	}

	node_json, err := json.MarshalIndent(node_metric, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprint(w, string(node_json))
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
func GetAppResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// get app Name
	appName := vars["appName"]

	var summary_cpu float64
	var memory_bundle []*docker.CgroupMemStat

	container_ids, pod_ips, err := GetContainerIDList("http://localhost", "8080", appName, "default")
	if err != nil {
		fmt.Println(err)
	}
	for _, container_id := range container_ids {
		fmt.Println(container_id, pod_ips)
		// calculation percentage of cpu usage
		container_cpu, _ := DockerCPUPercent(time.Duration(1) * time.Second, container_id)
		summary_cpu += container_cpu
		// memory usage 
		container_memory, _ := docker.CgroupMemDocker(container_id)
		memory_bundle = append(memory_bundle, container_memory)
	}

	// find the request per sec from haproxy-frontend
	res_front, err := http.Get("http://localhost:10001/v1/stats/frontends")
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
	res_back, err := http.Get("http://localhost:10001/v1/stats/backends")
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

	fmt.Println("rps: ", rps, ", rtime: ", rtime)
	// find the cpu avarage of application cpu usage
	average_cpu := summary_cpu/float64(len(container_ids))
	// create appliction object
	app_metric := AppMetric{
		App: appName,
		Cpu: average_cpu,
		Memory: memory_bundle,
		Request: rps_int,
		Response: rtime_int,
		Response2xx: res2xx_int,
		Response4xx: res4xx_int,
		Response5xx: res5xx_int,
	}

	app_json, err := json.MarshalIndent(app_metric, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprint(w, string(app_json))

}