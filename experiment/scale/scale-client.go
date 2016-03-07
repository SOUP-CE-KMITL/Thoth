package main

import (
	"github.com/google/cadvisor/client"
	info "github.com/google/cadvisor/info/v1"
)

func main() {

	sInfo, err := client.ContainerInfo("/docker/d9d3eb10179e6f93a...", &request)
}
