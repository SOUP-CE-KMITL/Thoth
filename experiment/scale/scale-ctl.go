package main

import (
	"fmt"
	"os/exec"
	"regexp"
)

func main() {
	fmt.Println("welcome to kubenetes scale-out experiment ... ")
	scale_num := "1"

	// check rc is running or not
	rc_stat, _ := exec.Command("kubectl", "get", "rc").Output()
	fmt.Println(string(rc_stat))
	re := regexp.MustCompile("nginx-controller")
	rc_exist := re.FindString(string(rc_stat))
	fmt.Println(string(rc_exist))
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

	fmt.Printf("replicas : ")
	fmt.Scanf("%s", &scale_num)
	cmd, _ := exec.Command("kubectl", "scale", "--replicas="+scale_num, "rc", "nginx-controller").Output()
	//cmd, _ := exec.Command("echo", "test").Output()
	fmt.Println(string(cmd))
}
