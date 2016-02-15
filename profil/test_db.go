package profil

//
//import (
//	"fmt"
//	"log"
//	"math/rand"
//	//"net/http"
//	//"os"
//	"github.com/influxdata/influxdb/client/v2"
//	"time"
//)
//
//var MyDB string = "thoth"
//var username string = "thoth"
//var password string = "thoth"
//
//func main2() {
//	// Make client
//	c, _ := client.NewHTTPClient(client.HTTPConfig{
//		Addr:     "http://localhost:8086",
//		Username: username,
//		Password: password,
//	})
//
//	// batch write
//	//	writePoints(c)
//	// query
//	res, err := queryDB(c, fmt.Sprintf("SELECT count(busy) FROM cpu_usage"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Print(res)
//}
//
//func writePoints(clnt client.Client) {
//	sampleSize := 1000
//	rand.Seed(42)
//
//	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
//		Database:  MyDB,
//		Precision: "us",
//	})
//
//	for i := 0; i < sampleSize; i++ {
//		regions := []string{"us-west1", "us-west2", "us-west3", "us-east1"}
//		tags := map[string]string{
//			"cpu":    "cpu-total",
//			"host":   fmt.Sprintf("host%d", rand.Intn(1000)),
//			"region": regions[rand.Intn(len(regions))],
//		}
//
//		idle := rand.Float64() * 100.0
//		fields := map[string]interface{}{
//			"idle": idle,
//			"busy": 100.0 - idle,
//		}
//
//		pt, _ := client.NewPoint("cpu_usage", tags, fields, time.Now())
//		bp.AddPoint(pt)
//	}
//
//	err := clnt.Write(bp)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
//	q := client.Query{
//		Command:  cmd,
//		Database: MyDB,
//	}
//	if response, err := clnt.Query(q); err == nil {
//		if response.Error() != nil {
//			return res, response.Error()
//		}
//		res = response.Results
//	} else {
//		return res, err
//	}
//	return res, nil
//}
