package main

import (
	"fmt"
	"github.com/white-pony/go-fann"
)

func main() {
	const numLayers = 3
	const desiredError = 0.00001
	const maxEpochs = 500000
	const epochsBetweenReports = 1000

	ann := fann.CreateStandard(numLayers, []uint32{2, 3, 1})
	ann.SetActivationFunctionHidden(fann.SIGMOID_SYMMETRIC)
	ann.SetActivationFunctionOutput(fann.SIGMOID_SYMMETRIC)
	ann.TrainOnFile("xor.data", maxEpochs, epochsBetweenReports, desiredError)

	fmt.Println("Testing network")

	test_data := fann.ReadTrainFromFile("xor.data")

	ann.ResetMSE()

	var i uint32
	for i = 0; i < test_data.Length(); i++ {
		fmt.Println(ann.Test(test_data.GetInput(i), test_data.GetOutput(i)))
	}

	fmt.Printf("MSE error on test data: %f\n", ann.GetMSE())
	fmt.Println("[0,0]", ann.Run([]fann.FannType{0, 0}))
	fmt.Println("[0,1]", ann.Run([]fann.FannType{0, 1}))
	fmt.Println("[1,0]", ann.Run([]fann.FannType{1, 0}))
	fmt.Println("[1,1]", ann.Run([]fann.FannType{1, 1}))
	ann.Save("xor_float.net")
	ann.Destroy()
}
