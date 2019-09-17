package main

import (
	"fmt"
	"log"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {

	tfModel, err := tf.LoadSavedModel("model/3", []string{"serve"}, nil)
	if err != nil {
		log.Panicf("Could not load model files into tensorflow with error: %v", err)
	}
	defer tfModel.Session.Close()
	ops := tfModel.Graph.Operations()
	for _, op := range ops {
		att, err := op.Attr("signature_def")
		if err == nil {
			fmt.Println(att, err)
		}
		fmt.Println(op.Name())
		// fmt.Println(att)
	}
	// op1 := tfModel.Graph.Operation("input").Operation(0)
	// fmt.Println(op1)
}
