package Dao

type WebhookMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					Context struct {
						From string `json:"from"`
						ID   string `json:"id"`
					} `json:"context"`
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Type      string `json:"type"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Sticker struct {
						MimeType string `json:"mime_type"`
						Sha256   string `json:"sha256"`
						ID       string `json:"id"`
					} `json:"sticker"`
					Document struct {
						Caption  string `json:"caption"`
						MimeType string `json:"mime_type"`
						Sha256   string `json:"sha256"`
						ID       string `json:"id"`
					} `json:"document"`
					Image struct {
						Caption  string `json:"caption"`
						MimeType string `json:"mime_type"`
						Sha256   string `json:"sha256"`
						ID       string `json:"id"`
					} `json:"image"`
					Video struct {
						Caption  string `json:"caption"`
						MimeType string `json:"mime_type"`
						Sha256   string `json:"sha256"`
						ID       string `json:"id"`
					} `json:"video"`
					Audio struct {
						Caption  string `json:"caption"`
						MimeType string `json:"mime_type"`
						Sha256   string `json:"sha256"`
						ID       string `json:"id"`
						Duration string `json:"duration"`
					}
					Errors struct {
						Code    int    `json:"code"`
						Details string `json:"details"`
						Title   string `json:"title"`
					} `json:"errors"`
					Button struct {
						Text    string `json:"text"`
						Payload string `json:"payload"`
					} `json:"button"`
					Interactive struct {
						ButtonReply struct {
							ID    string `json:"id"`
							Title string `json:"title"`
						} `json:"button_reply"`
						ListReply struct {
							ID          string `json:"id"`
							Title       string `json:"title"`
							Description string `json:"description"`
						} `json:"list_reply"`
						Type string `json:"type"`
					} `json:"interactive"`
					Referral struct {
						SourceURL    string `json:"source_url"`
						SourceID     string `json:"source_id"`
						SourceType   string `json:"source_type"`
						Headline     string `json:"headline"`
						Body         string `json:"body"`
						MediaType    string `json:"media_type"`
						ImageURL     string `json:"image_url"`
						VideoURL     string `json:"video_url"`
						ThumbnailURL string `json:"thumbnail_url"`
					} `json:"referral"`
					System struct {
						Body    string `json:"body"`
						NewWaID string `json:"new_wa_id"`
						Type    string `json:"type"`
					} `json:"system"`
				} `json:"messages"`
				Statuses []struct {
					ID           string `json:"id"`
					RecipientID  string `json:"recipient_id"`
					Status       string `json:"status"`
					Timestamp    string `json:"timestamp"`
					Conversation struct {
						ID                  string `json:"id"`
						ExpirationTimestamp string `json:"expiration_timestamp"`
						Origin              struct {
							Type string `json:"type"`
						} `json:"origin"`
					} `json:"conversation"`
					Pricing struct {
						PricingModel string `json:"pricing_model"`
						Billable     bool   `json:"billable"`
						Category     string `json:"category"`
					} `json:"pricing"`
					Errors []struct {
						Code  int    `json:"code"`
						Title string `json:"title"`
					} `json:"errors"`
				} `json:"statuses"`
				Errors struct {
					Code    int    `json:"code"`
					Details string `json:"details"`
					Title   string `json:"title"`
				} `json:"errors"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

type WebhookResponse struct {
	Entry  []Entry `json:"entry"`
	Object string  `json:"object"`
}

type Entry struct {
	Changes []Change `json:"changes"`
	ID      string   `json:"id"`
}

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Value struct {
	Contacts         []Contact `json:"contacts,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
	Statuses         []Status  `json:"statuses,omitempty"`
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WAID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From      string   `json:"from"`
	ID        string   `json:"id"`
	Text      Text     `json:"text"`
	Timestamp string   `json:"timestamp"`
	Type      string   `json:"type"`
	Reaction  Reaction `json:"reaction"`
}

type Text struct {
	Body string `json:"body"`
}

type Status struct {
	Conversation Conversation `json:"conversation"`
	ID           string       `json:"id"`
	Pricing      Pricing      `json:"pricing"`
	RecipientID  string       `json:"recipient_id"`
	Status       string       `json:"status"`
	Timestamp    string       `json:"timestamp"`
}

type Conversation struct {
	ID     string `json:"id"`
	Origin Origin `json:"origin"`
}

type Origin struct {
	Type string `json:"type"`
}

type Pricing struct {
	Billable     bool   `json:"billable"`
	Category     string `json:"category"`
	PricingModel string `json:"pricing_model"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Reaction struct {
	Emoji     string `json:"emoji"`
	MessageID string `json:"message_id"`
}
