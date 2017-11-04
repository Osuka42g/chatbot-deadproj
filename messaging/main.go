package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	accessToken      = ""
	facebookEndpoint = "https://graph.facebook.com/v2.6/me/messages?access_token=" + accessToken

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
		json.NewEncoder(w).Encode(standardResponse{"ok"})
		s := fbSenderInformation{}
		s.id, s.kind, s.payload = parseFBRequest(r)
		if s.kind != "invalid" {
			sendFBResponse(s)
		}
		return
	default:
		sendBadRequest(w, "Method not supported")
	}
}

func sendFBResponse(rs fbSenderInformation) {
	res := fbResponse{}
	res.Recipient.ID = rs.id
	res.Message.Text = rs.payload

	payload, _ := json.Marshal(res)
	req, err := http.NewRequest("POST", facebookEndpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func parseFBRequest(r *http.Request) (string, string, string) {
	fb := fbRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&fb)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	sender := fb.Entry[0].Messaging[0].Sender.SenderID
	kind := ""
	payload := ""
	if fb.Entry[0].Messaging[0].Message.Text != "" {
		kind = "text"
		payload = fb.Entry[0].Messaging[0].Message.Text
	} else if len(fb.Entry[0].Messaging[0].Message.Attachment) > 0 {
		kind = fb.Entry[0].Messaging[0].Message.Attachment[0].Type
		payload = fb.Entry[0].Messaging[0].Message.Attachment[0].Payload.URL
	} else {
		kind = "invalid"
	}
	return sender, kind, payload
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
	res := standardResponse{m}
	json.NewEncoder(w).Encode(res)
}
