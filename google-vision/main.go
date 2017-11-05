package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	vision "google.golang.org/api/vision/v1"
)

const developerKey = ""

type requestPayload struct {
	MediaURL string `json:"media_url"`
	Features string `json:"features"`
}

func main() {
	sendToGV()
	http.HandleFunc("/request", request)
	http.ListenAndServe(":8003", nil)

}

func request(w http.ResponseWriter, r *http.Request) {
	sendToGV()
}

func analyzeImageFromWeb(u string) {

}

func analyzeImage(i []byte) {

}

func saveImage(url string, output string) (int64, error) {
	img, _ := os.Create(output)
	resp, _ := http.Get(url)
	return io.Copy(img, resp.Body)
}

func sendToGV() {
	data, err := ioutil.ReadFile("./ferret.jpg")

	enc := base64.StdEncoding.EncodeToString(data)
	img := &vision.Image{Content: enc}

	feature := &vision.Feature{
		Type:       "LABEL_DETECTION",
		MaxResults: 10,
	}

	featurer := &vision.Feature{
		Type:       "LOGO_DETECTION",
		MaxResults: 10,
	}

	req := &vision.AnnotateImageRequest{
		Image:    img,
		Features: []*vision.Feature{feature, featurer},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	svc, err := vision.New(client)
	if err != nil {
		log.Fatal(err)
	}
	res, err := svc.Images.Annotate(batch).Do()
	if err != nil {
		log.Fatal(err)
	}

	body, err := json.Marshal(res.Responses[0].LabelAnnotations)
	fmt.Println(string(body))
}
