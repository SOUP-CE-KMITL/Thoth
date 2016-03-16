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
		// Get All RC name

		rcList, err := http.Get(thoth.KubeApi + "/api/v1/replicationcontrollers")
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(rcList.Body)
		rcList.Body.Close()
		if err != nil {
			panic(err)
		}

		var objRc interface{}
		if err := json.Unmarshal([]byte(body), &objRc); err != nil {
			panic(err)
		}

		// Extract RC Name
		RCArray := []thoth.RC{}
		RCLen := 0
		_RCLen := len(objRc.(map[string]interface{})["items"].([]interface{}))
		for i := 0; i < _RCLen; i++ {
			namespace := objRc.(map[string]interface{})["items"].([]interface{})[i].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
			if namespace != "default" {
				rc := thoth.RC{
					Name:      objRc.(map[string]interface{})["items"].([]interface{})[i].(map[string]interface{})["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["app"].(string),
					Namespace: namespace,
				}
				fmt.Println(rc.Namespace + "/" + rc.Name)
				RCArray = append(RCArray, rc)
				RCLen++
			}
		}
		//fmt.Println(RCArray)

		// Getting App resource usage
		for i := 0; i < RCLen; i++ {
			res := profil.GetAppResource(RCArray[i].Namespace, RCArray[i].Name)
			//	fmt.Println(res)
			replicas := profil.GetReplicas(RCArray[i].Namespace, RCArray[i].Name)
			// thoth.AppMetric{App,Cpu,Memory,Request,Response,Response2xx,Response4xx,Response5xx

			tags := map[string]string{
				"name": RCArray[i].Name,
			}
			
			fields := map[string]interface{}{
				// CPU
				// Memory
				"cpu":      res.Cpu,
				"memory":   res.Memory,
				"request":  res.Request * int64(30) / int64(replicas),
				"response": res.Response,
				"code2xx":  res.Response2xx,
				"code4xx":  res.Response4xx,
				"code5xx":  res.Response5xx,
				"replicas": replicas,
			}
			// fmt.Println(fields)
			if err := profil.WritePoints(c, MyDB, RCArray[i].Namespace, "s", tags, fields); err != nil {
				panic(err)
			}
			queryRes, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT count(response) FROM "+RCArray[i].Namespace))
			//queryRes, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT * FROM "+RCArray[i].Namespace))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(queryRes)
		}
		fmt.Println("Sleep")
		time.Sleep(10 * time.Second)
	}
	//-------------------------------------------------------------
}
