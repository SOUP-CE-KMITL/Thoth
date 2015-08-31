package main

import (
	"fmt"
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
	rm := memMet{memThreshold: metric_value}
	testScale(rm)

	// looping for check resource usage limits

}
