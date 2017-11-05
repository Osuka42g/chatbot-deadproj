package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	serverURL = "http://localhost:8082"
)

type msgRequest struct {
	id      string `json:"id"`
	kind    string `json:"kind"`
	payload string `json:"payload"`
}

func main() {
	http.HandleFunc("/middleware", request)
	http.ListenAndServe(":8002", nil)
}

func request(w http.ResponseWriter, r *http.Request) {

	s1 := "http://www.gcfa.com/assets/images/theme/home-ferret-trio-01.jpg"

	if !isValidURL(s1) {
		return
	}
	_, err := saveImage(s1)
	if err != nil {
		panic(err)
	}
}

func saveImage(URL string) (filepath string, err error) {
	p, _ := url.ParseRequestURI(URL)
	ext := path.Ext(strings.Split(p.RequestURI(), "?")[0]) // Get extension of the file, without downloading yet
	now := int(time.Now().Unix())
	filepath = "./downloads/" + strconv.Itoa(now) + ext

	createDownloadsDir()
	img, err := os.Create(filepath)
	resp, err := http.Get(URL)
	w, err := io.Copy(img, resp.Body)
	if err != nil {
		return
	}

	fmt.Println("Saved " + filepath + " " + strconv.Itoa(int(w)) + "bytes")
	return
}

func isValidURL(URL string) bool {
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return false
	}
	return true
}

func createDownloadsDir() {
	dir := "donwloads"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
}
