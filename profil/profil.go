package profil

import (
	"fmt"
	influx "github.com/influxdata/influxdb/client/v2"
	"strconv"
)

func GetProfilAvg(conn influx.Client, namespace, name, field, timeLength string) (float64, error) {
	query := fmt.Sprint("SELECT MEAN(" + field + ") FROM " + namespace + " WHERE app =~ /" + name + "/ AND time > now() - " + timeLength)
	dbResult, err := QueryDB(conn, query)
	if err != nil {
		panic(err)
		return -1.0, err
	}
	//	fmt.Println(dbResult)
	return strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][1]), 32)
}

func GetProfilLast(conn influx.Client, namespace, name, timeLength string) map[string]float64 {
	//mean(cpu),mean(memory),mean(request),mean(response),stddev(code5xx)
	query := fmt.Sprint("SELECT mean(cpu) as cpu,mean(memory) as memory,mean(rps) as rps,mean(rtime) as rtime,stddev(r2xx) as r2xx,stddev(r5xx) as r5xx FROM " + namespace + " WHERE app =~ /" + name + "/ AND time > now() - " + timeLength)
	dbResult, err := QueryDB(conn, query)
	if err != nil {
		panic(err)
		return nil
	}
	res := make(map[string]float64)
	res[fmt.Sprint(dbResult[0].Series[0].Columns[1])], _ = strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][1]), 32)
	res[fmt.Sprint(dbResult[0].Series[0].Columns[2])], _ = strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][2]), 32)
	res[fmt.Sprint(dbResult[0].Series[0].Columns[3])], _ = strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][3]), 32)
	res[fmt.Sprint(dbResult[0].Series[0].Columns[4])], _ = strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][4]), 32)
	res[fmt.Sprint(dbResult[0].Series[0].Columns[5])], _ = strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][5]), 32)
	res[fmt.Sprint(dbResult[0].Series[0].Columns[6])], _ = strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][6]), 32)
	return res
}

func WriteRPI(conn influx.Client, namespace, name string, request int64, replicas int) error {

	tags := map[string]string{
		"app": name,
	}

	fields := map[string]interface{}{
		"rps":      request,
		"replicas": replicas,
		"rpi":      request / int64(replicas), // TODO:remove?
	}
	fmt.Println(fields)
	if err := WritePoints(conn, namespace+"_rpi", "s", tags, fields); err != nil {
		return err
	}
	return nil
}

func GetAvgRPI(conn influx.Client, namespace, name string) (float64, error) {
	query := fmt.Sprint("SELECT MEAN(rpi) FROM " + namespace + "_rpi WHERE app =~ /" + name + "/ limit 20")
	dbResult, err := QueryDB(conn, query)
	if err != nil {
		panic(err)
		return -1.0, err
	}
	//	fmt.Println(dbResult)
	return strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][1]), 32)
}
