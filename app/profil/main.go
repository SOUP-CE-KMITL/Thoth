package main

import (
	"fmt"
	"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/SOUP-CE-KMITL/Thoth/profil"
	"github.com/influxdata/influxdb/client/v2"
	"log"
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
		runningPod := profil.GetAllRunningPod()
		RC := profil.GetUserRC()
		SVC := profil.GetUserSVC()
		RCLen := len(RC)

		statsMap := profil.GetHAProxyStats()
		// Getting App resource usage
		for i := 0; i < RCLen; i++ {
			appName := RC[i].Namespace + "/" + RC[i].Name
			value, exist := runningPod[appName]
			if exist && value {
				res := profil.GetAppResource(RC[i].Namespace, RC[i].Name)
				replicas, err := profil.GetReplicas(RC[i].Namespace, RC[i].Name)
				if err != nil {
					log.Println(err)
				}

				tags := map[string]string{
					"app": RC[i].Name,
				}
				svcSpec := SVC[appName]
				fmt.Println(svcSpec)
				//			if svcSpec != nil {
				fmt.Println(SVC[appName].Spec.Ports)
				nodePort := SVC[appName].Spec.Ports[0].NodePort
				vampPort := nodePort - 21000
				//			fmt.Println("VAMPStats ", vampStats.Spec.Ports[0].NodePort)
				fmt.Println("Port ", vampPort)
				//			fmt.Println("VAMPStats ", statsMap[vampPort])
				stats := statsMap[vampPort]
				fields := map[string]interface{}{
					"cpu":       res.Cpu,
					"memory":    res.Memory,
					"rps":       stats.Request,
					"rtime":     stats.Response,
					"r2xx":      stats.Response2xx,
					"r4xx":      stats.Response4xx,
					"r5xx":      stats.Response5xx,
					"r5xxroute": stats.Response5xxRoute,
					"replicas":  replicas,
				}
				fmt.Println(fields)
				if err := profil.WritePoints(c, RC[i].Namespace, "s", tags, fields); err != nil {
					panic(err)
				}
				//			}
			}
		}
		fmt.Println("Sleep\n")
		time.Sleep(10 * time.Second)
	}
	//-------------------------------------------------------------
}
