package main

import (
	"encoding/json"
	//	"runtime"
	//	"strconv"
	"os/exec"
	//"strings"
	"fmt"
	"io/ioutil"
	"net/http"
)

var thothApiUrl string = "http://localhost:8182"

func main() {
	// Get all RC
	jsonRc := GetApps()
	fmt.Print(string(jsonRc))
	var objRc interface{}
	err := json.Unmarshal((jsonRc), &objRc)
	if err != nil {
		panic(err)
	}
	// Extract RC Name
	rcName := objRc.(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
	fmt.Print(rcName)

	//http://localhost:8182/app/<rc-name>/metrics
	res, err := http.Get(thothApiUrl + "/app/" + rcName + "/metrics")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	} else {
		var resJson map[string]interface{}
		if err := json.Unmarshal(body, &resJson); err != nil {
			panic(err)
		}
		fmt.Print(resJson["Response"].(float64))
	}
}

func GetApps() []byte {
	res, err := exec.Command("kubectl", "get", "rc", "-o", "json").Output()
	if err != nil {
		panic(err)
	}
	return res
}
