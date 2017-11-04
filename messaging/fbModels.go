package main

type fbResponse struct {
	Message string `json:"message"`
}

type fbRequest struct {
	Object string           `json:"object"`
	Entry  []fbRequestEntry `json:"entry"`
}

type fbRequestEntry struct {
	Messaging []fbRequestMessaging `json:"messaging"`
}

type fbRequestMessaging struct {
	Sender struct {
		SenderID string `json:"id"`
	} `json:"sender"`
	Message struct {
		Text       string                `json:"text"`
		Attachment []fbRequestAttachment `json:"attachments"`
	} `json:"message"`
}

type fbRequestAttachment struct {
	Type    string `json:"type"`
	Payload struct {
		URL string `json:"url"`
	} `json:"payload"`
}
