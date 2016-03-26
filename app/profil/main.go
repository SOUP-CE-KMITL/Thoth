package main

import (
	"encoding/json"
	"fmt"
	"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/SOUP-CE-KMITL/Thoth/profil"
	"github.com/influxdata/influxdb/client/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var MyDB string = "thoth"
var username string = "thoth"
var password string = "thoth"

func main() {
	// Connect InfluxDB
	c, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr:     thoth.InfluxdbApi,
		Username: username,
		Password: password,
	})

	//-------------------------------------------------------------
	for {
		// Get All RC&SVC name
		RC := profil.GetUserRC()
		SVC := profil.GetUserSVC()
		fmt.Println(SVC)
		RCLen := len(RC)

		// Getting App resource usage
		for i := 0; i < RCLen; i++ {
			fmt.Println(RC[i].Namespace + "/" + RC[i].Name)
			res := profil.GetAppResource(RC[i].Namespace, RC[i].Name)
			//	fmt.Println(res)
			replicas, err := profil.GetReplicas(RC[i].Namespace, RC[i].Name)
			if err != nil {
				//	panic(err)
				log.Println(err)
			}
			// thoth.AppMetric{App,Cpu,Memory,Request,Response,Response2xx,Response4xx,Response5xx

			tags := map[string]string{
				"app": RC[i].Name,
			}

			fields := map[string]interface{}{
				// CPU
				// Memory
				"cpu":    res.Cpu,
				"memory": res.Memory, /*
					"request":  res.Request,
					"response": res.Response,
					"code2xx":  res.Response2xx,
					"code4xx":  res.Response4xx,
					"code5xx":  res.Response5xx,
				*/
				"replicas": replicas,
			}
			fmt.Println(fields)
			if err := profil.WritePoints(c, RC[i].Namespace, "s", tags, fields); err != nil {
				panic(err)
			}
		}
		fmt.Println("Sleep\n")
		time.Sleep(10 * time.Second)
	}
	//-------------------------------------------------------------
}

func getHAProxyStats() {
	// find the request per sec from haproxy-frontend
	resFront, err := http.Get(thoth.VampApi + "/v1/stats/frontends")
	if err != nil {
		panic(err)
	}
	bodyFront, err := ioutil.ReadAll(resFront.Body)
	resFront.Body.Close()
	if err != nil {
		panic(err)
	}
	//var rps uint64
	var objectFront thoth.VampFrontend
	err = json.Unmarshal([]byte(bodyFront), &objectFront)
	if err != nil {
		fmt.Println(err)
	}

	//find resonse time from haproxy-backends
	resBack, err := http.Get(thoth.VampApi + "/v1/stats/backends")
	if err != nil {
		panic(err)
	}
	bodyBack, err := ioutil.ReadAll(resBack.Body)
	resBack.Body.Close()
	if err != nil {
		panic(err)
	}
	var objectBack thoth.VampBackend
	err = json.Unmarshal([]byte(bodyBack), &objectBack)
	if err != nil {
		fmt.Println(err)
	}
}
func getHAProxyStatsApp(nodePort string) (response, request, res2xx, res4xx, res5xx int64) {
	// find the request per sec from haproxy-frontend
	resFront, err := http.Get(thoth.VampApi + "/v1/stats/frontends")
	if err != nil {
		panic(err)
	}
	bodyFront, err := ioutil.ReadAll(resFront.Body)
	resFront.Body.Close()
	if err != nil {
		panic(err)
	}
	//var rps uint64
	var objectFront []map[string]interface{}
	err = json.Unmarshal([]byte(bodyFront), &objectFront)
	strRps := objectFront[0]["req_rate"].(string)
	request, _ = strconv.ParseInt(strRps, 10, 64)
	if err == nil {
	} else {
		fmt.Println(err)
	}

	//find resonse time from haproxy-backends
	resBack, err := http.Get(thoth.VampApi + "/v1/stats/backends")
	if err != nil {
		panic(err)
	}
	bodyBack, err := ioutil.ReadAll(resBack.Body)
	resBack.Body.Close()
	if err != nil {
		panic(err)
	}
	var objectBack []map[string]interface{}
	err = json.Unmarshal([]byte(bodyBack), &objectBack)
	strResponse := objectBack[0]["rtime"].(string)
	strRes2xx := objectBack[0]["hrsp_2xx"].(string)
	strRes4xx := objectBack[0]["hrsp_4xx"].(string)
	strRes5xx := objectBack[0]["hrsp_5xx"].(string)
	response, _ = strconv.ParseInt(strResponse, 10, 64)
	res2xx, _ = strconv.ParseInt(strRes2xx, 10, 64)
	res4xx, _ = strconv.ParseInt(strRes4xx, 10, 64)
	res5xx, _ = strconv.ParseInt(strRes5xx, 10, 64)
	if err == nil {
	} else {
		fmt.Println(err)
	}
	return
}
