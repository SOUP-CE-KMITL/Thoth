package main

import (
	"fmt"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		fmt.Println("error : ", e)
	}
}

func main() {
	// TODO : resource usage
	dat, err := ioutil.ReadFile("/sys/fs/cgroup/memory/docker/6e6bfe5c5fe1b6de929501c00da9f838cf0a2591d839308362e51a7f32ecba3d/memory.usage_in_bytes")
	check(err)
	fmt.Print(string(dat))

}
