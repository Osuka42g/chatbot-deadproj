package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	accessToken      = ""
	facebookEndpoint = "https://graph.facebook.com/v2.6/me/messages?access_token=" + accessToken

	verificationToken  = "AwesomeYouMadeAGreatJob"
	middlewareEndpoint = "https://e0f8f652.ngrok.io/middleware" // We still don't have, but we will
)

func main() {
	http.HandleFunc("/messenger", routeMessage)
	http.ListenAndServe(":8001", nil)
}

func routeMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET": // The only GET method facebook will send us, is for the verification challenge.
		verifyFacebookChallenge(w, r)
	case "POST":
		handleFBPostRequest(w, r)
	default:
		respondBadRequest(w, "Method not supported")
	}
}

func handleFBPostRequest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(standardResponse{"ok"})
	ss := fbSenderInformation{}
	ss.Id, ss.Kind, ss.Payload = parseFBRequest(r)

	if ss.Kind == "invalid" {
		return
	}
	sendFBPayload(composeFBTyping(ss, true))
	time.Sleep(2 * time.Second) // Sleep 2 seconds to be more natural

	payload, _ := json.Marshal(ss)
	mw, err := fetchFromMiddleware(payload)
	if err != nil {
		panic(err)
	}
	ss.Payload = mw

	err = sendFBPayload(composeFBMessage(ss))
	if err != nil {
		panic(err)
	} else {
		sendFBPayload(composeFBTyping(ss, false))
	}
}

func composeFBMessage(rs fbSenderInformation) []byte {
	res := fbSimpleText{}
	res.Recipient.ID = rs.Id
	res.Message.Text = rs.Payload
	payload, _ := json.Marshal(res)
	return payload
}

func composeFBTyping(rs fbSenderInformation, mode bool) []byte {
	res := fbTyping{}
	res.Recipient.ID = rs.Id
	res.SenderAction = "typing_off"
	if mode {
		res.SenderAction = "typing_on"
	}
	payload, _ := json.Marshal(res)
	return payload
}

func sendFBPayload(p []byte) error {
	req, err := http.NewRequest("POST", facebookEndpoint, bytes.NewBuffer(p))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func fetchFromMiddleware(p []byte) (response string, err error) {
	response = ""
	err = nil
	req, err := http.NewRequest("POST", middlewareEndpoint, bytes.NewBuffer(p))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	msg := fbSenderInformation{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&msg)

	response = msg.Payload
	return
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
	kind := "invalid"
	payload := ""
	message := fb.Entry[0].Messaging[0].Message

	if message.Text != "" {
		kind = "text"
		payload = message.Text
	} else if len(message.Attachment) > 0 {
		kind = message.Attachment[0].Type
		payload = message.Attachment[0].Payload.URL
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
	respondBadRequest(w, "Invalid verification token")
}

func respondBadRequest(w http.ResponseWriter, m string) {
	w.WriteHeader(http.StatusBadRequest)
	res := standardResponse{m}
	json.NewEncoder(w).Encode(res)
}
