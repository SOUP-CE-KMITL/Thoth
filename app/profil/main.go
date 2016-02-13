package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	//"html"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/SOUP-CE-KMITL/Thoth/profil"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

var kubeApi string = "http://localhost:8080"
var influxdbApi string = "http://localhost:8086"
var vampApi string = "http://localhost:10001"

var MyDB string = "thoth"
var username string = "thoth"
var password string = "thoth"

type RC struct {
	Namespace string // Namespace = User
	Name      string
}

func main() {
	// Connect InfluxDB
	c, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxdbApi,
		Username: username,
		Password: password,
	})

	//-------------------------------------------------------------
	for {
		// Get All RC name
		jsonRc := GetPods()
		var objRc interface{}
		err := json.Unmarshal([]byte(jsonRc), &objRc)
		if err != nil {
			panic(err)
		}
		// Extract RC Name
		RCArray := []RC{}
		RCLen := 0
		_RCLen := len(objRc.(map[string]interface{})["items"].([]interface{}))
		for i := 0; i < _RCLen; i++ {
			namespace := objRc.(map[string]interface{})["items"].([]interface{})[i].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
			if namespace != "default" {
				rc := RC{
					Name:      objRc.(map[string]interface{})["items"].([]interface{})[i].(map[string]interface{})["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["app"].(string),
					Namespace: namespace,
				}
				fmt.Println(rc.Namespace + "/" + rc.Name)
				RCArray = append(RCArray, rc)
				RCLen++
			}
		}
		//fmt.Println(RCArray)

		// Getting App resource usage
		for i := 0; i < RCLen; i++ {
			res := GetAppResource(RCArray[i].Namespace, RCArray[i].Name)
			//	fmt.Println(res)

			// thoth.AppMetric{App,Cpu,Memory,Request,Response,Response2xx,Response4xx,Response5xx

			tags := map[string]string{
				"name": RCArray[i].Name,
			}

			fields := map[string]interface{}{
				// CPU
				// Memory
				"cpu":      res.Cpu,
				"memory":   res.Memory,
				"request":  res.Request,
				"response": res.Response,
				"code2xx":  res.Response2xx,
				"code4xx":  res.Response4xx,
				"code5xx":  res.Response5xx,
			}
			fmt.Println(fields)
			if err := profil.WritePoints(c, MyDB, RCArray[i].Namespace, "s", tags, fields); err != nil {
				panic(err)
			}
			queryRes, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT count(response) FROM "+RCArray[i].Namespace))
			//queryRes, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT * FROM "+RCArray[i].Namespace))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(queryRes)
		}
		fmt.Println("Sleep")
		time.Sleep(10 * time.Second)
	}
	//-------------------------------------------------------------
}

// list every node
func GetNodes(w http.ResponseWriter, r *http.Request) {
	// to do need to read api and port of api server from configuration file
	res, err := http.Get(kubeApi + "/api/v1/nodes")
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
	res, err := http.Get(kubeApi + "/api/v1/nodes/" + nodesName)
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

// list all pods
func GetPods() string {
	// to do need to read api and port of api server from configuration file
	res, err := http.Get(kubeApi + "/api/v1/pods")
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

// list specific pod details
func GetPod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// node name from user.
	podName := vars["podName"]
	fmt.Fprint(w, string(podName))
	// to do need to read api and port of api server from configuration file
	// TODO: change namespace to flexible.
	var dat map[string]interface{}
	res, err := http.Get(kubeApi + "/api/v1/namespaces/default/pods/" + podName)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	pretty_body, err := json.MarshalIndent(dat, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(pretty_body))
}

// post handler for scale pod by pod name

func GetApp(w http.ResponseWriter, r *http.Request) {

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

func GetApps(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// node name from user.
	namespace := vars["namespace"]

	res, err := exec.Command("kubectl", "get", "rc", "-o", "json", "--namespace="+namespace).Output()
	fmt.Println("namespace = " + namespace)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(res))
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
		fmt.Println("Busy: ", t2Busy-t1Busy, ", All: ", t2All-t1All)
		fmt.Println("idle time 1: ", t1Busy-t1All, ", idle time 2: ", t2Busy-t2All)
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

	container_ids, pod_ips, err := GetContainerIDList(kubeApi, name, namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, container_id := range container_ids {
		fmt.Println(container_id, pod_ips)
		// calculation percentage of cpu usage
		container_cpu, _ := DockerCPUPercent(time.Duration(1)*time.Second, container_id)
		summary_cpu += container_cpu
		// memory usage
		container_memory, _ := docker.CgroupMemDocker(container_id)
		memory_bundle = append(memory_bundle, container_memory)
	}

	podNum := len(pod_ips)

	// find the request per sec from haproxy-frontend
	res_front, err := http.Get(vampApi + "/v1/stats/frontends")
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
	res_back, err := http.Get(vampApi + "/v1/stats/backends")
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
	//app_json, err := json.MarshalIndent(app_metric, "", "\t")
	//if err != nil {
	//	fmt.Println("error:", err)
	//}
	//fmt.Fprint(w, string(app_json))

}
