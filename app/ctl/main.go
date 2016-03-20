package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/white-pony/go-fann"
	"strings"
	//"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/SOUP-CE-KMITL/Thoth/profil"
	//	"io/ioutil"
	"log"
	//	"math"
	"encoding/csv"
	"github.com/goml/gobrain"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var username string = "thoth"
var password string = "thoth"
var MyDB string = "thoth"

//var influxdbApi string = "127.0.0.1:8086"

func main() {
	for {

		// Get all user RC
		RCArray := profil.GetUserRC()
		RCLen := len(RCArray)

		// Getting App Metric
		for i := 0; i < RCLen; i++ {
			replicas, err := profil.GetReplicas(RCArray[i].Namespace, RCArray[i].Name)

			//				res, err := thoth.ScaleOutViaCli(1, RCArray[i].Namespace, RCArray[i].Name)
			if err != nil {
				panic(err)
			}
			fmt.Println(replicas)

			// Check & Label & Save WPI

		}
		// -----Prediction-----
		// Normalize
		// Run (Predict)
		// Label
		runFann()
		//-----------
		fmt.Println("Sleep")
		time.Sleep(1 * time.Minute)
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
