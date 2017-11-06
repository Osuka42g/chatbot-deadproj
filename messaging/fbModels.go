package main

type standardResponse struct {
	Message string `json:"message"`
}

type fbRequest struct {
	Entry []fbRequestEntry `json:"entry"`
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

type fbSenderInformation struct {
	Id      string `json:"id"`
	Kind    string `json:"kind"` // Alias for `type`, reserved world in go
	Payload string `json:"payload"`
}

type fbSimpleText struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

type fbTyping struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	SenderAction string `json:"sender_action"`
}
