package Dao

type MessageResponse struct {
	Object string `json:"object"`
	Entry  []struct {
		Id      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberId      string `json:"phone_number_id"`
				}
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					}
					WaId string `json:"wa_id"`
				}
				Messages []struct {
					From      string `json:"from"`
					Id        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					}
					Type string `json:"type"`
				}
			}
			Field string `json:"field"`
		}
	}
}
