// TODO : should change to api package, it's should call by main core
// package api
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
)

// pod metrics
type Pod struct {
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	Cpu       uint64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
	DiskUsage uint64 `json:"disk_usage"`
	Bandwidth uint64 `json:"bandwidth"`
}

// array of pods type
type Pods []Pod

// node metrics
type Node struct {
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	Cpu       uint64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
	DiskUsage uint64 `json:"disk_usage"`
}

// array of node type
type Nodes []Node

// TODO : need to change to init ,this server it should run by core main.
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/metrics/", Metrics)
	router.HandleFunc("/mock", nodeTestMock)

	log.Fatal(http.ListenAndServe(":8181", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Metrics")
}

// TODO : remove
// test mocks
func nodeTestMock(w http.ResponseWriter, r *http.Request) {
	nodes := Nodes{
		Node{Name: "node1", Ip: "192.168.1.2", Cpu: 5000, Memory: 3000, DiskUsage: 1000},
		Node{Name: "node2", Ip: "192.168.1.4", Cpu: 5000, Memory: 3000, DiskUsage: 1000},
	}

	json.NewEncoder(w).Encode(nodes)
}
