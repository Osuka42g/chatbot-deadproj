package main

import (
	"encoding/json"
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
	Kind    string `json:"kind"`
	Payload string `json:"payload"`
}

func main() {
	http.HandleFunc("/middleware", request)
	http.ListenAndServe(":8002", nil)
}

func request(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		fmt.Fprint(w, "Method not supported")
		return
	}

	ss := msgRequest{}
	ss = parseMsgRequest(r)

	switch ss.Kind {
	case "image":
		// res, err := examineImage(ss.Payload)
		ss.Payload = "mm that's not a ferret!"
	case "text":
		if "help" == ss.Payload {
			ss.Payload = "Just send me the ferrets!"
		} else {
			ss.Payload = randomInvalid()
		}
	}

	ss.Kind = "text"
	json.NewEncoder(w).Encode(ss)
}

func randomInvalid() string {
	messages := []string{
		"the ferrets!",
		"show me the ferrets!!",
		"aaaarrrrrrrrgg",
	}
	return messages[0]
}

func examineImage(i string) (result string, err error) {
	result = ""
	err = nil
	if !isValidURL(i) {
		return
	}
	_, err = saveImage(i)
	if err != nil {
		panic(err)
	}

	return
}

func parseMsgRequest(r *http.Request) msgRequest {
	msg := msgRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	return msg
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
	dir := "downloads"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
}
