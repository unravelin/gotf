package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type predReqData struct {
	Inputs [][][][]float32 `json:"inputs"`
}

func NewCustomTargeter(name, url string, data [][]byte) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}
		tgt.Method = "POST"
		tgt.URL = url
		h := http.Header{}
		b := data[rand.Intn(len(data))]
		switch name {
		case "compiled-zipped":
			var buf bytes.Buffer
			zw, _ := gzip.NewWriterLevel(&buf, -1)
			_, err := zw.Write(b)
			if err != nil {
				panic(err)
			}
			if err := zw.Close(); err != nil {
				panic(err)
			}
			h.Set("Content-Encoding", "gzip")
			tgt.Body = buf.Bytes()
		default:
			tgt.Body = b
			h.Set("Content-Type", "application/json")
		}
		tgt.Header = h
		return nil
	}
}

func genInput(images int) [][][][]float32 {
	a := make([][][][]float32, images)
	for i := range a {
		a[i] = make([][][]float32, 224)
		for j := range a[i] {
			a[i][j] = make([][]float32, 224)
			for k := range a[i][j] {
				a[i][j][k] = make([]float32, 3)
				for l := range a[i][j][k] {
					a[i][j][k][l] = rand.Float32()
				}
			}
		}
	}
	return a
}

func main() {
	var dataSlice [][]byte

	images := 1

	freqLst := []int{1, 2, 5, 10, 15, 25, 50}
	duration := 3 * time.Second
	randoms := 10
	urls := map[string]string{
		// "compile": "http://localhost:8502/v1/models/resnet_model:predict",
		"serving": "http://localhost:8501/v1/models/resnet_model:predict",
	}
	fmt.Printf("Backend\tImages\tMean     \tP99     \tP95     \tMax     \tSuccess      \tThroughput\n")
	// for _, images := range imagesLst {
	for _, freq := range freqLst {
		// Make some fake data
		dataSlice = make([][]byte, randoms)
		for i := 0; i < randoms; i++ {
			r := predReqData{Inputs: genInput(images)}

			b, err := json.Marshal(r)

			// _, err := zw.Write(b)
			if err != nil {
				panic(err)
			}
			dataSlice[i] = b
		}

		rate := vegeta.Rate{Freq: freq, Per: time.Second}

		metrics := make(map[string]*vegeta.Metrics)
		// success = true

		for name, url := range urls {
			metrics[name] = &vegeta.Metrics{}
			targeter := NewCustomTargeter(name, url, dataSlice)
			attacker := vegeta.NewAttacker()

			for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
				metrics[name].Add(res)
				// if res.Success != 1 {

				// }
			}
			metrics[name].Close()

			fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%f\t%f\n", name, freq, metrics[name].Latencies.Mean, metrics[name].Latencies.P99, metrics[name].Latencies.P95, metrics[name].Latencies.Max, metrics[name].Success, metrics[name].Throughput)
			// if metrics[name].Success < 1 {
			// 	fmt.Println("No longer 100%% success, stopping this model")
			// 	// success = false
			// 	delete(urls, name)
			// 	// break
			// }
			time.Sleep(1 * time.Second)
		}

	}
}

// r := predReqData{Inputs: genInput()}
// b, err := json.Marshal(r)

// if err != nil {
// 	fmt.Printf("Error: %s", err)
// 	panic(err)
// }

// req, err := http.NewRequest("POST", "http://localhost:8501/v1/models/resnet_model:predict", bytes.NewBuffer(b))
// if err != nil {
// 	panic(err)
// }
// req.Header.Set("Content-Type", "application/json")

// client := &http.Client{}
// resp, err := client.Do(req)
// if err != nil {
// 	panic(err)
// }
// defer resp.Body.Close()

// fmt.Println("response Status:", resp.Status)
// panic("here")
