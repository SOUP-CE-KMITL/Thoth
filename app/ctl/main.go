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
	"os/exec"
	"strconv"
	"time"
)

var username string = "thoth"
var password string = "thoth"
var MyDB string = "thoth"

//var influxdbApi string = "127.0.0.1:8086"

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
		/*
			jsonRc := profil.GetPods()
			var objRc interface{}
			err := json.Unmarshal([]byte(jsonRc), &objRc)
			if err != nil {
				panic(err)
			}
		*/

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

			/*Name,Tags,Columns,Values,Err*/
			//fmt.Println(len(qres10min[0].Series[0].Values))
			//fmt.Println(len(qres1hr[0].Series[0].Values))

			qresRep, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT LAST(replicas) FROM "+RCArray[i].Namespace))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(qresRep)
			// Check for prevent error bcuz no data
			if len(qresRep[0].Series) != 0 {
				replicas, err := strconv.ParseInt(fmt.Sprint(qresRep[0].Series[0].Values[0][1]), 10, 32)
				if err != nil {
					panic(err)
				}
				// SCALE 1
				//if replicas == 1 {
				qres1hr, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT MEAN(response) FROM "+RCArray[i].Namespace+" WHERE time > now() - 1d")) //+" WHERE app="RCArray[i].Name+))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("1D", qres1hr)

				qres10min, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT MEAN(response) FROM "+RCArray[i].Namespace+" WHERE time > now() - 10m"))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("10Min", qres10min)
				res1hr, err := strconv.ParseFloat(fmt.Sprint(qres1hr[0].Series[0].Values[0][1]), 32)
				if err == nil {
					fmt.Println(res1hr)
					if res1hr != 0 {
						res10min, err := strconv.ParseFloat(fmt.Sprint(qres10min[0].Series[0].Values[0][1]), 32)
						if err == nil {
							fmt.Println(res10min)
							if err == nil {
								if res10min != 0 {
									/* qresSpread, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT spread(response) FROM "+RCArray[i].Namespace+" WHERE time > now() - 1d"))
									if err != nil {
										log.Fatal(err)
									}
									fmt.Println("ResSpread"+qresSpread)
									resSpread, err := strconv.ParseFloat(fmt.Sprint(qresSpread[0].Series[0].Values[0][1]), 32)*/
									//									if res10min-res1hr > resSpread {
									if res10min > res1hr {
										// Delay each scale 10 min

										qRepSpread, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT spread(replicas) FROM "+RCArray[i].Namespace+" WHERE time > now() - 10m"))
										if err != nil {
											log.Fatal(err)
										}
										fmt.Println(qRepSpread)
										repSpread, err := strconv.ParseFloat(fmt.Sprint(qRepSpread[0].Series[0].Values[0][1]), 32)
										if repSpread < 1 {
											fmt.Println("Scale +1")
											res, err := scaleOutViaCli(int(replicas)+1, RCArray[i].Namespace, RCArray[i].Name)
											if err != nil {
												panic(err)
											}
											fmt.Println(res)
										}
									} else if replicas > 1 {
										fmt.Println("Scale -1")
										res, err := scaleOutViaCli(int(replicas)-1, RCArray[i].Namespace, RCArray[i].Name)
										if err != nil {
											panic(err)
										}
										fmt.Println(res)
									}
								}
							}
						}
					}
				}
			}
			//else {
			//	// SCALE 2 - Condition ( avg10min(req) >
			//	qReq1d, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT MEAN(request) FROM "+RCArray[i].Namespace+" WHERE time > now() - 1d"))
			//	if err != nil {
			//		log.Fatal(err)
			//	}
			//	fmt.Println(qReq1d)
			//	qReq1d, err := profil.QueryDB(c, MyDB, fmt.Sprint("SELECT MEAN(request) FROM "+RCArray[i].Namespace+" WHERE time > now() - 1d"))
			//	if err != nil {
			//		log.Fatal(err)
			//	}
			//	fmt.Println(qReq1d)
			//}
			// If 5xx then Scalee
			//			}
		}
		fmt.Println("Sleep")
		time.Sleep(5 * time.Minute)
	}
}

/**
* scale-out replicas via cli
**/
func scaleOutViaCli(replicas int, namespace, name string) (string, error) {
	var err error
	var cmd []byte
	fmt.Println("s_n : ", replicas, "rc : ", name)
	if cmd, err = exec.Command("kubectl", "scale", "--replicas="+strconv.Itoa(replicas), "rc", name, "--namespace", namespace).Output(); err != nil {
		fmt.Println(err)
	}
	return string(cmd), err
}
