package bluejoy

type PushbulletHTTReply struct {
	Accounts      []interface{} `json:"accounts"`
	Blocks        []interface{} `json:"blocks"`
	Channels      []interface{} `json:"channels"`
	Chats         []interface{} `json:"chats"`
	Clients       []interface{} `json:"clients"`
	Contacts      []interface{} `json:"contacts"`
	Devices       []interface{} `json:"devices"`
	Grants        []interface{} `json:"grants"`
	Pushes        []Push        `json:"pushes"`
	Profiles      []interface{} `json:"profiles"`
	Subscriptions []interface{} `json:"subscriptions"`
	Texts         []interface{} `json:"texts"`
	Cursor        string        `json:"cursor"`
}

type Push struct {
	Active                  bool     `json:"active"`
	Iden                    string   `json:"iden"`
	Created                 float64  `json:"created"`
	Modified                float64  `json:"modified"`
	Type                    string   `json:"type"`
	Dismissed               bool     `json:"dismissed"`
	GUID                    string   `json:"guid"`
	Direction               string   `json:"direction"`
	SenderIden              string   `json:"sender_iden"`
	SenderEmail             string   `json:"sender_email"`
	SenderEmailNormalized   string   `json:"sender_email_normalized"`
	SenderName              string   `json:"sender_name"`
	ReceiverIden            string   `json:"receiver_iden"`
	ReceiverEmail           string   `json:"receiver_email"`
	ReceiverEmailNormalized string   `json:"receiver_email_normalized"`
	SourceDeviceIden        string   `json:"source_device_iden"`
	AwakeAppGuids           []string `json:"awake_app_guids"`
	Title                   string   `json:"title"`
	URL                     string   `json:"url"`
}
