package main

import (
	"fmt"
	"http"
	"io/ioutil"
)

func main() {

}

func GetContainerIDList(url string, port string, rc_name string, namespace string) (string, error) {
	res, err := http.Get(url + ":" + port + "/api/v1/watch/namespaces/" + namespace + "/pods")
	if err != nil {
		fmt.Println("Can't connect to cadvisor")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return 0, err
	} else {
		// TODO: remove
		fmt.Println("response : " + body)
		// json handler type
		var res_obj map[string]interface{}
		if err := json.Unmarshal(body, &res_obj); err != nil {
			return 0, err
		}
		//var res_time = res_obj["stats"].([]interface{})[len(res_obj["stats"].([]interface{}))-1].(map[string]interface{})["timestamp"].(string)
		/*
			var container_id string
			if container_id, err = GetContainerIdByReplication(rc_name); err != nil {
				return 0, err
			}
			return container_id, nil // success retrieve resource usage info
		*/
		// this is for testing
		//return uint64(res_obj["stats"].([]interface{})[0].(map[string]interface{})["memory"].(map[string]interface{})["usage"].(float64)), res_time, nil // success retrieve resource usage info
	}
	// TODO : remove this it's for dummy return
	return "error", error.New("dummy error")
}
