package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	contianers, ips, _ := GetContainerIDList("http://localhost", "8080", "goweb-controller", "default")
	fmt.Println(contianers, ips)
}

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
