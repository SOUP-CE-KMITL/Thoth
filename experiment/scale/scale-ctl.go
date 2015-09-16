package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type metricsController interface {
	getThreshold() uint64
	setThreshold(metrics uint64) bool
	checkMeetThreshold(current uint64) bool
}

type latencyMet struct {
	latencyThreshold uint64
}

type memMet struct {
	// TODO : type should be uint64
	memThreshold uint64
}

type cpuMet struct {
	// TODO : cpu type should be float64
	cpuThreshold uint64
}

// implement PID Controller
// Latency implemented
func (l latencyMet) setThreshold(metrics uint64) bool {
	l.latencyThreshold = metrics
	// TODO : return error if have somethings wrong
	return true
}
func (l latencyMet) getThreshold() uint64 {
	return l.latencyThreshold
}
func (l latencyMet) checkMeetThreshold(current uint64) bool {
	return l.latencyThreshold == current
}

// Memory Resource implemented
func (rm memMet) setThreshold(metrics uint64) bool {
	rm.memThreshold = metrics
	// TODO : return error if have somethings wrong
	return true
}
func (rm memMet) getThreshold() uint64 {
	return rm.memThreshold
}
func (rm memMet) checkMeetThreshold(current uint64) bool {
	return rm.memThreshold == current
}

// CPU Resource implemented
func (rc cpuMet) getThreshold() uint64 {
	return rc.cpuThreshold
}
func (rc cpuMet) setThreshold(metrics uint64) bool {
	rc.cpuThreshold = metrics
	// TODO : return error if have somethings wrong
	return true
}
func (rc cpuMet) checkMeetThreshold(current uint64) bool {
	return rc.cpuThreshold == current
}

// TODO : remove this func when finish
func testScale(m metricsController) {
	// TODO : delete below line it's just for test
	fmt.Println(m)
	fmt.Println("Threshold : ", m.getThreshold())
}

/**
/	Reader function
**/

func check(e error) {
	if e != nil {
		fmt.Println("ERROR : ", e)
	}
}

/**
/	json encoder for get current resource information
**/
func GetResourceUsage(resource_type int, res_obj map[string]interface{}) (uint64, error) {
	switch resource_type {
	case 2:
		return uint64(res_obj["stats"].([]interface{})[0].(map[string]interface{})["cpu"].(map[string]interface{})["usage"].(map[string]interface{})["total"].(float64)), nil
	case 3:
		return uint64(res_obj["stats"].([]interface{})[0].(map[string]interface{})["memory"].(map[string]interface{})["usage"].(float64)), nil
		// TODO : latency is not implement for now.
	default:
		return 0, errors.New("invaild choice, Please select between 1-4") // TODO : need to catch this error
	}
}

/**
/	Instantiate function for create new metrics controller obj.
**/
func MetricFactory(metric_type int, threshold uint64) (metricsController, error) {
	switch metric_type {
	case 2:
		return cpuMet{cpuThreshold: threshold}, nil
	case 3:
		return memMet{memThreshold: threshold}, nil
	default:
		return nil, errors.New("invaild choice, Please select between 1-4")
	}
}

/**
/	api call
**/
func GetCurrentResource(url string, metric_type int) (uint64, string, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Can't connect to cadvisor")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return 0, "", err
	} else {
		// json handler type
		var res_obj map[string]interface{}
		if err := json.Unmarshal(body, &res_obj); err != nil {
			return 0, "", err
		}
		var res_time = res_obj["stats"].([]interface{})[len(res_obj["stats"].([]interface{}))-1].(map[string]interface{})["timestamp"].(string)
		var resource_usage uint64
		if resource_usage, err = GetResourceUsage(metric_type, res_obj); err != nil {
			return 0, "", err
		}
		return resource_usage, res_time, nil // success retrieve resource usage info
		// this is for testing
		//return uint64(res_obj["stats"].([]interface{})[0].(map[string]interface{})["memory"].(map[string]interface{})["usage"].(float64)), res_time, nil // success retrieve resource usage info
	}
}

/**
* 	Get current resource by using cgroup
**/
func GetCurrentResourceCgroup(container_id string, metric_type int) (uint64, error) {
	// TODO : Read Latency from HA Proxy

	// file path prefix
	var path = "/sys/fs/cgroup/memory/docker/" + container_id + "/"

	if metric_type == 2 {
		// read cpu usage
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
* scale-out replicas via cli
**/
func scaleOutViaCli(scale_num int, rc_name string) (string, error) {
	var err error
	var cmd []byte
	fmt.Println("s_n : ", scale_num, "rc : ", rc_name)
	if cmd, err = exec.Command("kubectl", "scale", "--replicas="+strconv.Itoa(scale_num), "rc", rc_name).Output(); err != nil {
		fmt.Println(err)
	}
	return string(cmd), err
}

func main() {

	fmt.Println("welcome to kubenetes scale-out experiment ... ")
	scale_num := 1
	var rc_name = "default"
	fmt.Println("please specific your replicaiton controller name ")
	fmt.Scanf("%s", &rc_name)

	// check rc is running or not
	rc_stat, _ := exec.Command("kubectl", "get", "rc").Output()
	fmt.Println(string(rc_stat))
	re := regexp.MustCompile(rc_name)
	rc_exist := re.FindString(string(rc_stat))
	fmt.Println(string(rc_exist))

	//  create new nginx rc if not created
	if rc_exist == "" {
		var yml_file string
		fmt.Println("you don't have any nginx rc waiting for a min to create..")
		fmt.Println("please fill path of rc yml files")
		fmt.Scanf("%s", &yml_file)
		// create new nginx rc from file
		_, err := exec.Command("kubectl", "create", "-f", yml_file).Output()
		if err != nil {
			fmt.Println("failed to create rc")
		}
	} else {
		fmt.Println("you already have rc ... ")
	}

	// TODO : pod creation is take time for a while so, we can't suddenly check after created.
	// double check pod is already created or not.
	/*pod, _ := exec.Command("kubectl", "get", "pod").Output()
	regex_pod := regexp.MustCompile(rc_name)
	pod_exist := regex_pod.FindString(string(pod))
	// TODO : remove this line
	fmt.Println("pod : ", string(pod), "pod_exists : ", pod_exist, " , rc_name = ", rc_name)
	if pod_exist != rc_name {
		err := errors.New("pod is not created !!")
		panic(err)
	}*/

	// call scale command from kubectl
	fmt.Printf("replicas : ")
	fmt.Scanf("%d", &scale_num)
	if cmd, err := scaleOutViaCli(scale_num, rc_name); err != nil {
		fmt.Println(err, cmd)
	}

	// set Threshold
	metric_type := 1
	fmt.Printf("\n Which metrics do you want to using ? \n")
	fmt.Printf("1. Latency \n")
	fmt.Printf("2. Cpu Usage \n")
	fmt.Printf("3. Memory Usage \n")
	fmt.Scanf("%d", &metric_type)
	// TODO : catch invalid choice value
	var metric_value uint64
	fmt.Printf("value = ")
	fmt.Scanf("%d", &metric_value)
	// TODO : catch invalid value

	// create new metrics object
	var metric metricsController
	var err error
	if metric, err = MetricFactory(metric_type, metric_value); err != nil {
		fmt.Println(err)
	}

	//  TEST : obj type
	fmt.Println("Threshold : ", metric.getThreshold())

	// ==== looping for check resource usage limits ====
	//var container_url = "http://localhost:4194/api/v1.0/containers/docker/7b1fa7a7d61b903a18aa14c230407b7b0302aaa2bc241fec1e3ded3f73cd96a2"
	var current_resource uint64
	// var res_time string
	// for test
	count_scale_time := 0
	// intial array of container_id and pod_ip

	container_ids, pod_ips, err := GetContainerIDList("http://localhost", "8080", rc_name, "default")
	if err != nil {
		fmt.Println(err)
	}
	// TODO : remove this print it's for dummy
	fmt.Println(pod_ips, container_ids)
	count_round_robin := 0
	for {
		// round robin get resource
		count_round_robin += 1
		count_round_robin = count_round_robin % len(container_ids)
		current_resource, err = GetCurrentResourceCgroup(container_ids[count_round_robin], metric_type)
		check(err)
		fmt.Println("Current usage at :", count_round_robin, " of ", len(container_ids), current_resource)
		if current_resource >= metric_value {
			if count_scale_time%60 == 0 {
				fmt.Println("reached threshold try to scale-out ...")
				scale_num++
				_, err := scaleOutViaCli(scale_num, rc_name)
				if err != nil {
					// got new array of container_ids
					container_ids, pod_ips, err = GetContainerIDList("http://localhost", "8080", rc_name, "default")
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
		fmt.Println("current resource at ", " : ", current_resource)
		time.Sleep(1000 * time.Millisecond)

	}
}
