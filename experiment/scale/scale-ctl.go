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
)

type metricsController interface {
	getThreshold() int
	setThreshold(metrics int) bool
	checkMeetThreshold(current int) bool
}

type latencyMet struct {
	latencyThreshold int
}

type memMet struct {
	// TODO : type should be uint64
	memThreshold int
}

type cpuMet struct {
	// TODO : cpu type should be float64
	cpuThreshold int
}

// json handler
type container struct {
}

// implement PID Controller
// Latency implemented
func (l latencyMet) setThreshold(metrics int) bool {
	l.latencyThreshold = metrics
	// TODO : return error if have somethings wrong
	return true
}
func (l latencyMet) getThreshold() int {
	return l.latencyThreshold
}
func (l latencyMet) checkMeetThreshold(current int) bool {
	return l.latencyThreshold == current
}

// Memory Resource implemented
func (rm memMet) setThreshold(metrics int) bool {
	rm.memThreshold = metrics
	// TODO : return error if have somethings wrong
	return true
}
func (rm memMet) getThreshold() int {
	return rm.memThreshold
}
func (rm memMet) checkMeetThreshold(current int) bool {
	return rm.memThreshold == current
}

// CPU Resource implemented
func (rc cpuMet) getThreshold() int {
	return rc.cpuThreshold
}
func (rc cpuMet) setThreshold(metrics int) bool {
	rc.cpuThreshold = metrics
	// TODO : return error if have somethings wrong
	return true
}
func (rc cpuMet) checkMeetThreshold(current int) bool {
	return rc.cpuThreshold == current
}

func testScale(m metricsController) {
	// TODO : delete below line it's just for test
	fmt.Println(m)
	fmt.Println("Threshold : ", m.getThreshold())
}

func GetResourceUsage(resource_type int, res_obj map[string]interface{}) (int, error) {
	switch resource_type {
	case 2:
		return int(res_obj["stats"].([]interface{})[0].(map[string]interface{})["cpu"].(map[string]interface{})["usage"].(map[string]interface{})["total"].(float64)), nil
	case 3:
		return int(res_obj["stats"].([]interface{})[0].(map[string]interface{})["memory"].(map[string]interface{})["usage"].(float64)), nil
		// TODO : latency is not implement for now.
	default:
		return 0, errors.New("invaild choice, Please select between 1-4") // TODO : need to catch this error
	}
}

func MetricFactory(metric_type int, threshold int) (metricsController, error) {
	switch metric_type {
	case 2:
		return cpuMet{cpuThreshold: threshold}, nil
	case 3:
		return memMet{memThreshold: threshold}, nil
	default:
		return nil, errors.New("invaild choice, Please select between 1-4")
	}
}

func main() {
	// intial error for catch any error.
	var err errors

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

	metric_value := 0
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
	if metric, err = MetricFactory(metric_type, metric_value); err != nil {
		panic(err)
	}

	//  TEST : obj type
	fmt.Println("Threshold : ", metric.getThreshold())

	// ==== looping for check resource usage limits ====
	// TODO : now I'm hardcode container name ,so it's need to get url wihtout hardcode
	res, err := http.Get("http://localhost:4194/api/v1.0/containers/docker/8f91896b6d2350d8623ce536ef8c986dd524bc1787813093132ee7d3050a6bf2")
	if err != nil {
		fmt.Println("Can't connect to cadvisor")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println("Can't read response body")
		log.Fatal(err)
	} else {
		// json handler type
		var res_obj map[string]interface{}
		if err := json.Unmarshal(body, &res_obj); err != nil {
			panic(err)
		}

		var resource_usage int
		if resource_usage, err = GetResourceUsage(metric_type, res_obj); err != nil {
			panic(err)
		}

	}
}
