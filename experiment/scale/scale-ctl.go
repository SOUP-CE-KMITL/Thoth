package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
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
	json encoder for get current resource information
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
	Instantiate function for create new metrics controller obj.
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
	api call
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
	}
}

func main() {

	fmt.Println("welcome to kubenetes scale-out experiment ... ")
	scale_num := "1"

	// check rc is running or not
	rc_stat, _ := exec.Command("kubectl", "get", "rc").Output()
	fmt.Println(string(rc_stat))
	re := regexp.MustCompile("nginx-controller")
	rc_exist := re.FindString(string(rc_stat))
	fmt.Println(string(rc_exist))

	//  create new nginx rc if not created
	if rc_exist == "" {
		yml_file := ""
		fmt.Println("you don't have any nginx rc waiting for a min to create..")
		fmt.Println("please fill path of nginx yml files")
		fmt.Scanf("%s", &yml_file)
		// create new nginx rc from file
		_, err := exec.Command("kubectl", "run", "-f", yml_file).Output()
		if err != nil {
			// TODO: need to real checking is container is already created via this command or not
			fmt.Println("successful to created nginx rc")
		} else {
			fmt.Println("failed to create nginx rc")
		}
	} else {
		fmt.Println("you already have nginx rc ... ")
	}

	// call scale command from kubectl
	fmt.Printf("replicas : ")
	fmt.Scanf("%s", &scale_num)
	cmd, _ := exec.Command("kubectl", "scale", "--replicas="+scale_num, "rc", "nginx-controller").Output()
	fmt.Println(string(cmd))

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

	// Sample object for testing
	/*
		rm := memMet{memThreshold: metric_value}
		testScale(rm)
	*/

	// create new metrics object
	var metric metricsController
	var err error
	if metric, err = MetricFactory(metric_type, metric_value); err != nil {
		panic(err)
	}

	//  TEST : obj type
	fmt.Println("Threshold : ", metric.getThreshold())

	// ==== looping for check resource usage limits ====
	// TODO : now I'm hardcode container name ,so it's need to get url wihtout hardcode
	var container_url = "http://localhost:4194/api/v1.0/containers/docker/339a7596578fcaa1d831e0f28c5bdc7f15e56675ca01120c2f223df928a4e5df"
	var current_resource uint64
	var res_time string

	for true {
		current_resource, res_time, err = GetCurrentResource(container_url, metric_type)
		fmt.Println("current resource at ", res_time, " : ", current_resource)
		time.Sleep(1000 * time.Millisecond)
	}
}
