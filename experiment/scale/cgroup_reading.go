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
	dat, err := ioutil.ReadFile("/sys/fs/cgroup/memory/docker/1145bc121145d04dbff622d4613c409751ceba6704f4f80be08e80e1cf596599/memory.usage_in_bytes")
	check(err)
	fmt.Print(string(dat))

}
