// TODO : need to change to api
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	// to do need to read api and port of api server from configuration file
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
	send_obj, err := json.Marshal(object)
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
	res, err := http.Get("http://localhost:8080/api/v1/nodes")
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
	// TODO: pods
	res, err := http.Get("http://localhost:8080/api/v1/pods/" + podName)
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

	jsonReq, err := json.Marshal(objReq)
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
