package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const facebookVerifyToken = "AwesomeYouMadeAGreatJob"

type facebookResponse struct {
	value string
}

func main() {
	http.HandleFunc("/messenger", GetMessage)
	http.ListenAndServe(":8001", nil)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" && r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		res := facebookResponse{"Invalid"}
		json.NewEncoder(w).Encode(res)
		return
	}

	// The only GET method facebook will send us, is for the verification challenge.
	if r.Method == "GET" {
		q := r.URL.Query()
		if q["hub.mode"][0] == "subscribe" && q["hub.verify_token"][0] == facebookVerifyToken {
			fmt.Fprintf(w, q["hub.challenge"][0])
		}
		return
	}

}
