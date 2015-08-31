package main

import (
	//	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Heelllo")
	dat, err := ioutil.ReadFile("/home/extraterrestrial/webserv-nginx.log2")
	var f []map[string]interface{}
	err = json.Unmarshal([]byte(dat), &f)

	if err != nil {
		fmt.Println(err)
	}

	fo, _ := os.Create("/mnt/nginx-parse.log")
	defer fo.Close()

	for _, v := range f {
		if v["stream"] == "stdout" {
			//fmt.Println(v["log"])
			_, err = fo.WriteString(v["log"].(string))
			check(err)
			fo.Sync()
		}
	}

}
