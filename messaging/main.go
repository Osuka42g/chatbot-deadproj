package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	accessToken        = ""
	verificationToken  = "AwesomeYouMadeAGreatJob"
	middlewareEndpoint = "" // We still don't have, but we will
)

func main() {
	http.HandleFunc("/messenger", getMessage)
	http.ListenAndServe(":8001", nil)
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	// The only GET method facebook will send us, is for the verification challenge.
	case "GET":
		verifyFacebookChallenge(w, r)
	case "POST":
		json.NewEncoder(w).Encode(fbResponse{"ok"})
		fb := fbRequest{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&fb)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		fmt.Println(fb)
	default:
		sendBadRequest(w, "Method not supported")
	}
}

func sendfbResponse() {

}

func verifyFacebookChallenge(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if len(q) != 3 {
	} else if q["hub.mode"][0] == "subscribe" && q["hub.verify_token"][0] == verificationToken {
		fmt.Fprintf(w, q["hub.challenge"][0])
		return
	}
	sendBadRequest(w, "Invalid verification token")
}

func sendBadRequest(w http.ResponseWriter, m string) {
	w.WriteHeader(http.StatusBadRequest)
	res := fbResponse{m}
	json.NewEncoder(w).Encode(res)
}
