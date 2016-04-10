package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/SOUP-CE-KMITL/Thoth/profil"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/white-pony/go-fann"
	"strings"
	//	"io/ioutil"
	"encoding/csv"
	"github.com/goml/gobrain"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var username string = "thoth"
var password string = "thoth"

type fannObj struct {
	rcName string
	input  []fann.FannType
	output []fann.FannType
}

func z_score(avg int64, sd float64, val int64) fann.FannType {
	if sd != 0 {
		return fann.FannType((float64(val) - float64(avg)) / sd)
	} else {
		return fann.FannType(0)
	}
}

var ann *fann.Ann

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			fmt.Println(sig)
			ann.Save("fann.dat")
			os.Exit(0)
		}
	}()
	// Connect InfluxDB
	influxDB, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     thoth.InfluxdbApi,
		Username: username,
		Password: password,
	})

	if err != nil {
		panic(err)
	}

	// FANN Init
	if _, err := os.Stat("fann.dat"); err == nil {
		// path/to/whatever exists
		fmt.Println("Load fann.dat")
		ann = fann.CreateFromFile("fann.dat")
	} else {
		fmt.Println("Init FANN")
		const num_layers = 3
		const num_neurons_hidden = 10
		const desired_error = 0.001

		ann = fann.CreateStandard(num_layers, []uint32{5, num_neurons_hidden, 1})
		ann.SetActivationFunctionHidden(fann.SIGMOID_SYMMETRIC)
		ann.SetActivationFunctionOutput(fann.SIGMOID_SYMMETRIC)
	}
	// TODO:it's need to calculate this for all RC
	avg_day := profil.GetProfilLast(influxDB, "thoth", "eight-puzzle", "1d")
	sd_day := profil.GetProfilStdLast(influxDB, "thoth", "eight-puzzle", "1d")

	for {
		// Get all user RCLen
		RC := profil.GetUserRC()
		RCLen := len(RC)

		f_obj := make([]fannObj, RCLen)

		// Getting App Metric
		for i := 0; i < RCLen; i++ {
			// Get 15 min avarage of all attribute for feed as input.
			resUsage15min := profil.GetProfilLast(influxDB, RC[i].Namespace, RC[i].Name, "15min")
			fmt.Println("count = ", i)
			f_obj[i].rcName = RC[i].Name
			// TODO : calculate z value for each attribute
			f_obj[i].input = []fann.FannType{z_score(avg_day["cpu"], sd_day["cpu"], resUsage15min["cpu"]),
				z_score(avg_day["memory"], sd_day["memory"], resUsage15min["memory"]),
				z_score(avg_day["rps"], sd_day["rps"], resUsage15min["rps"]),
				z_score(avg_day["rtime"], sd_day["rtime"], resUsage15min["rtime"]),
				z_score(avg_day["r2xx"], sd_day["r2xx"], resUsage15min["r2xx"]),
				z_score(avg_day["r5xx"], sd_day["r5xx"], resUsage15min["r5xx"]),
				z_score(avg_day["replicas"], sd_day["replicas"], resUsage15min["replicas"])}
			f_obj[i].output = []fann.FannType{0}

			replicas, err := profil.GetReplicas(RC[i].Namespace, RC[i].Name)

			if err != nil {
				panic(err)
			}
			fmt.Println(replicas)

			// Check Resposne time & Label & Save WPI
			var responseDay, response10Min float64
			if responseDay, err = profil.GetProfilAvg(influxDB, RC[i].Namespace, RC[i].Name, "rtime", "1d"); err != nil {
				panic(err)
				log.Println(err)
			}
			fmt.Println("resDays : ", responseDay)

			if response10Min, err = profil.GetProfilAvg(influxDB, RC[i].Namespace, RC[i].Name, "rtime", "5m"); err != nil {
				fmt.Println("res10min : ", response10Min)
				panic(err)
				log.Println(err)
			}
			// Floor
			responseDay = math.Floor(responseDay)
			response10Min = math.Floor(response10Min)
			fmt.Println("D", responseDay, " 10M", response10Min)
			metrics := profil.GetAppResource(RC[i].Namespace, RC[i].Name)
			var cpu10Min float64
			if cpu10Min, err = profil.GetProfilAvg(influxDB, RC[i].Namespace, RC[i].Name, "cpu", "5m"); err != nil {
				panic(err)
				log.Println(err)
			}
			fmt.Println("CPU ", cpu10Min)
			if cpu10Min > 70 {
				fmt.Println("Response check")
				if response10Min > responseDay { // TODO:Need to check WPI too
					// Save WPI
					fmt.Println("Scale+1")
					f_obj[i].output = []fann.FannType{1}
					if err := profil.WriteRPI(influxDB, RC[i].Namespace, RC[i].Name, metrics.Request, replicas); err != nil {
						panic(err)
						log.Println(err)
					}
					// Scale +1
					// TODO: Limit
					if replicas < 10 {
						if _, err := thoth.ScaleOutViaCli(replicas+1, RC[i].Namespace, RC[i].Name); err != nil {
							panic(err)
						}
					}
				}
			} else if replicas > 1 {
				// = rpi/replicas
				var rpiMax float64
				if rpiMax, err = profil.GetAvgRPI(influxDB, RC[i].Namespace, RC[i].Name); err != nil {
					rpiMax = -1
					// TODO:Handler
					//panic(err)
				}
				fmt.Println("WPI", rpiMax)
				if rpiMax > 0 {
					minReplicas := int(metrics.Request / int64(rpiMax)) // TODO: Ceil?

					if minReplicas < replicas {
						// Scale -1
						fmt.Println("Scale-1")
						f_obj[i].output = []fann.FannType{-1}
						if _, err := thoth.ScaleOutViaCli(replicas-1, RC[i].Namespace, RC[i].Name); err != nil {
							panic(err)
						}
					}
				}
			}

		}

		// -----Prediction-----
		// Normalize
		// Run (Predict)
		// Label
		fmt.Println("============================ RUNNING FANN ============================")
		runFann(f_obj)
		//-----------
		fmt.Println("Sleep TODO:Change to 5 Minnn")
		time.Sleep(1 * time.Minute)
	}
}

func runFann(fObj []fannObj) {

	// iterate to train data
	for i := range fObj {
		ann.Train(fObj[i].input, fObj[i].output)
		fmt.Println("Train Data ", fObj)
	}

	//ann.TrainOnData(train_data, 2000, 500, desired_error)
	fmt.Printf("MSE : %f\n", ann.GetMSE())
}

func annGoBrain() {
	// Open the file.
	f, _ := os.Open("file.csv")
	r := csv.NewReader(f)
	trainSet := [][][]float64{}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		record := [][]float64{
			{parse(rec[5]), parse(rec[6]), parse(rec[7]), parse(rec[8]), parse(rec[9])}, {parse(rec[10])},
		}
		fmt.Println(strings.Join(rec[5:], " "))
		fmt.Println(record)
		trainSet = append(trainSet, record)
	}
	//-----------
	rand.Seed(0)

	// instantiate the Feed Forward
	ff := &gobrain.FeedForward{}

	// initialize the Neural Network;
	// inputs, hidden nodes and output.
	ff.Init(5, 10, 1)

	// train the network using the XOR patterns,1000 epochs,learning rate 0.6,momentum factor 0.4,receive reports about error
	ff.Train(trainSet, 1000, 0.6, 0.4, true)

	//Test

	//	ff.Test([][][]float64{{{40, 70, 2, 1965, 13}, {1}}})
}

func parse(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
