package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client/v2"
	_ "log"
	"net/url"
	_ "os"
	"time"
)

// influxDB variable
const (
	MyDB     = "thoth"
	username = "thoth_user"
	password = "thoth_secret"
)

func main() {
	// Male client
	u, _ := url.Parse("http://localhost:8086")
	c := client.NewClient(client.Config{
		URL:      u,
		Username: username,
		Password: password,
	})

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})

	// Create point and add to batch
	// cpu is tags
	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}

	// cpu usage is measurement
	pt := client.NewPoint("cpu_usage", tags, fields, time.Now())
	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)

	res, _ := queryDB(c, "SELECT * FROM cpu_usage")
	fmt.Println(res)

}

func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}
