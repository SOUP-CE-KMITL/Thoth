package main

import (
	//"bufio"
	"encoding/csv"
	"fmt"
	"github.com/goml/gobrain"
	"io"
	"math/rand"
	"os"

	"log"
	"strconv"
	//	"strings"
)

func parse(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func main() {
	// Open the file.
	f, _ := os.Open("file.csv")
	//5,6,7,8,9,(10)
	// Create a new Scanner for the file.
	//	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.

	r := csv.NewReader(f)
	/*
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
	*/
	trainSet := [][][]float64{}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(record)
		record := [][]float64{
			{parse(rec[5]), parse(rec[6]), parse(rec[7]), parse(rec[8]), parse(rec[9])}, {parse(rec[10])},
		}
		//		fmt.Println(record)
		trainSet = append(trainSet, record)
	}
	fmt.Println(trainSet)
	//-----------
	rand.Seed(0)

	// create the XOR representation patter to train the network
	// req,
	/*
		patterns := [][][]float64{
			// 0
			{{2}, {0}},
			{{3}, {0}},
			{{1000}, {0}},
			{{800}, {0}},
			{{500}, {0}},
			{{600}, {0}},
			{{50}, {0}},
			{{100}, {0}},

			// +1
			{{1260}, {1}},
			{{1295}, {1}},
			{{1399}, {1}},
			{{1399}, {1}},
			{{1399}, {1}},
			{{1399}, {1}},
			{{1399}, {1}},
			{{1399}, {1}},
			{{1260}, {1}},
			{{2000}, {1}},
			{{1979}, {1}},
			{{2027}, {1}},
			{{1983}, {1}},
			{{2000}, {1}},
		}
	*/
	// instantiate the Feed Forward
	ff := &gobrain.FeedForward{}

	// initialize the Neural Network;
	// the networks structure will contain:
	// 2 inputs, 2 hidden nodes and 1 output.
	ff.Init(5, 5, 1)

	// train the network using the XOR patterns
	// the training will run for 1000 epochs
	// the learning rate is set to 0.6 and the momentum factor to 0.4
	// use true in the last parameter to receive reports about the learning error
	ff.Train(trainSet, 1000, 0.6, 0.4, true)

	//Test

	ff.Test([][][]float64{{{40, 70, 2, 1965, 13}, {1}}})
}
