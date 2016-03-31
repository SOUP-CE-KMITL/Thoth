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
		RC := profil.GetUserRC()
		SVC := profil.GetUserSVC()
		fmt.Println(SVC)
		RCLen := len(RC)

		statsMap := profil.GetHAProxyStats()
		fmt.Println(statsMap)
		// Getting App resource usage
		for i := 0; i < RCLen; i++ {
			fmt.Println(RC[i].Namespace + "/" + RC[i].Name)
			res := profil.GetAppResource(RC[i].Namespace, RC[i].Name)
			replicas, err := profil.GetReplicas(RC[i].Namespace, RC[i].Name)
			if err != nil {
				//	panic(err)
				log.Println(err)
			}

			tags := map[string]string{
				"app": RC[i].Name,
			}

			key := RC[i].Namespace + "/" + RC[i].Name
			svcSpec := SVC[key]
			fmt.Println(svcSpec)
			//			if svcSpec != nil {
			fmt.Println(SVC[key].Spec.Ports)
			nodePort := SVC[key].Spec.Ports[0].NodePort
			vampPort := nodePort - 21000
			//			fmt.Println("VAMPStats ", vampStats.Spec.Ports[0].NodePort)
			fmt.Println("Port ", vampPort)
			//			fmt.Println("VAMPStats ", statsMap[vampPort])
			stats := statsMap[vampPort]
			fields := map[string]interface{}{
				// CPU
				// Memory
				"cpu":          res.Cpu,
				"memory":       res.Memory,
				"request":      stats.Request,
				"response":     stats.Response,
				"code2xx":      stats.Response2xx,
				"code4xx":      stats.Response4xx,
				"code5xx":      stats.Response5xx,
				"code5xxroute": stats.Response5xxRoute,

				"replicas": replicas,
			}
			fmt.Println(fields)
			if err := profil.WritePoints(c, RC[i].Namespace, "s", tags, fields); err != nil {
				panic(err)
			}
			//			}
		}
		fmt.Println("Sleep\n")
		time.Sleep(10 * time.Second)
	}
	//-------------------------------------------------------------
}
