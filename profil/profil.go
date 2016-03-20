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

func WriteWPI(conn influx.Client, namespace, name string, request int64, replicas int) error {

	tags := map[string]string{
		"app": name,
	}

	fields := map[string]interface{}{
		"request":  request,
		"replicas": replicas,
		"wpi":      request / int64(replicas), // TODO:remove?
	}
	fmt.Println(fields)
	if err := WritePoints(conn, namespace+"_wpi", "s", tags, fields); err != nil {
		return err
	}
	return nil
}

func GetAvgWPI(conn influx.Client, namespace, name, timeLength string) (float64, error) {
	query := fmt.Sprint("SELECT MEAN(wpi) FROM " + namespace + "_wpi WHERE app =~ /" + name + "/ AND time > now() - " + timeLength)
	dbResult, err := QueryDB(conn, query)
	if err != nil {
		panic(err)
		return -1.0, err
	}
	//	fmt.Println(dbResult)
	return strconv.ParseFloat(fmt.Sprint(dbResult[0].Series[0].Values[0][1]), 32)
}
