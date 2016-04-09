package main

import (
	"encoding/csv"
	"fmt"
	"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/SOUP-CE-KMITL/Thoth/learn"
	"github.com/SOUP-CE-KMITL/Thoth/profil"
	"github.com/goml/gobrain"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/white-pony/go-fann"
	"io"
	"log"
	"os/signal"
	"strings"
	//"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var username string = "thoth"
var password string = "thoth"

func main() {
	agent := learn.QLearn{Gamma: 0.3, Epsilon: 0.3}
	agent.Init()
	if err := agent.Load("ql.da"); err != nil {
		fmt.Println("Load Fail", err)
	}
	agent.Epsilon = 0.3

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			fmt.Println(sig)
			agent.Save("ql.da")
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

	firstRun := true
	lastState := learn.State{}
	lastAction := 0
	for {
		// Get all user RC
		RC := profil.GetUserRC()
		RCLen := len(RC)

		// Getting App Metric
		for i := 0; i < RCLen; i++ {
			replicas, err := profil.GetReplicas(RC[i].Namespace, RC[i].Name)
			if err != nil {
				panic(err)
			}

			// TODO : Chane time interval
			res := profil.GetProfilLast(influxDB, RC[i].Namespace, RC[i].Name, "1m")
			fmt.Println(res)

			if !firstRun {
				// Reward Last state
				agent.Reward(lastState, lastAction, res)
			}
			firstRun = false

			lastState = agent.CurrentState
			agent.SetCurrentState(res["cpu"], res["memory"], res["rps"], res["rtime"], res["r5xx"], replicas)
			action := agent.ChooseAction()

			if action+replicas != 0 {
				if _, err := thoth.ScaleOutViaCli(replicas+action, RC[i].Namespace, RC[i].Name); err != nil {
					fmt.Println(err)
				}
			}
			lastAction = action
			fmt.Println(agent)
		}
		//-----------
		fmt.Println("Sleep TODO:Change to 5 Min")
		time.Sleep(60 * time.Second)
	}
}

func runFann() {
	const num_layers = 3
	const num_neurons_hidden = 10
	const desired_error = 0.001

	train_data := fann.ReadTrainFromFile("file.csv")
	//	test_data := fann.ReadTrainFromFile("../../datasets/robot.test")

	var momentum float32
	//	for momentum = 0.0; momentum < 0.7; momentum += 0.1 {
	fmt.Printf("============= momentum = %f =============\n", momentum)

	ann := fann.CreateStandard(num_layers, []uint32{train_data.GetNumInput(), num_neurons_hidden, train_data.GetNumOutput()})
	/*
		ann.SetTrainingAlgorithm(fann.TRAIN_INCREMENTAL)
		ann.SetLearningMomentum(momentum)
	*/
	ann.SetActivationFunctionHidden(fann.SIGMOID_SYMMETRIC)
	ann.SetActivationFunctionOutput(fann.SIGMOID_SYMMETRIC)
	ann.TrainOnData(train_data, 2000, 500, desired_error)

	fmt.Printf("MSE error on train data: %f\n", ann.TestData(train_data))
	//	fmt.Printf("MSE error on test data : %f\n", ann.TestData(test_data))

	ann.Destroy()
	//	}

	train_data.Destroy()
	//test_data.Destroy()
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
