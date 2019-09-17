package main

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {

	tfModel, err := tf.LoadSavedModel("model", []string{"serve"}, nil)
	if err != nil {
		log.Panicf("Could not load model files into tensorflow with error: %v", err)
	}
	defer tfModel.Session.Close()

	http.HandleFunc("/v1/models/resnet_model:predict", func(w http.ResponseWriter, r *http.Request) {
		predict(w, r, tfModel)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

var zippers = sync.Pool{New: func() interface{} {
	return gzip.NewWriter(nil)
}}

func predict(w http.ResponseWriter, r *http.Request, tfModel *tf.SavedModel) {

	var predReqData struct {
		Inputs [][][][]float32
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&predReqData); err != nil {
		w.WriteHeader(400)
		log.Printf("Error parsing request body: %v", err)
		return
	}

	tensor, err := tf.NewTensor(predReqData.Inputs)

	if err != nil {
		w.WriteHeader(400)
		log.Printf("Error creating tensors: %v", err)
		return
	}

	s, err := calculate(tfModel, tensor)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error predicting using tf model: %v", err)
		return
	}

	prediction, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling tf output to json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(prediction)
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func calculate(tfModel *tf.SavedModel, input *tf.Tensor) (result [][]float32, err error) {

	run, err := tfModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			tfModel.Graph.Operation("input_tensor").Output(0): input,
		},
		[]tf.Output{
			tfModel.Graph.Operation("softmax_tensor").Output(0),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return run[0].Value().([][]float32), nil
}
