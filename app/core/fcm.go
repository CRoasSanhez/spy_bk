package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/revel/revel"
)

// FCM Server constants
const (
	// FCM Server
	// FCMServerURL is base url to send push notification.
	FCMServerURL = "https://fcm.googleapis.com/fcm/send"

	// Dev keys
	DevFCMServerKey = "AIzaSyCz5XgBG9KMXDkrHPzIYrZFYyHxtkfRm50"
	DevMaxTTL       = 86400

	// Prod keys
	FCMServerKey = "AIzaSyDPPRi3U5Rd9mH9Pe36twVDgfsdrvmkKZw"
	MAXTTL       = 604800

	// Notification const
	PriorityHigh   = "high"
	PriorityNormal = "normal"
)

// FCMNotification struct for notification message
type FCMNotification struct {
	Body  string `json:"body,omitempty"`
	Title string `json:"title,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Badge string `json:"badge,omitempty"`
	Sound string `json:"sound,omitempty"`
}

// FCMClient stores the key and the Message (FcmMsg)
type FCMClient struct {
	Message FCMMessage
}

// FCMMessage represents fcm request message
type FCMMessage struct {
	To               string          `json:"to,omitempty"`
	RegistrationIds  []string        `json:"registration_ids,omitempty"`
	CollapseKey      string          `json:"collapse_key,omitempty"`
	Priority         string          `json:"priority,omitempty"`
	Data             interface{}     `json:"data,omitempty"`
	Notification     FCMNotification `json:"notification,omitempty"`
	ContentAvailable bool            `json:"content_available,omitempty"`
	TimeToLive       int             `json:"time_to_live,omitempty"`
}

// FcmResponseStatus represents fcm response message - (tokens and topics)
type FcmResponseStatus struct {
	Ok         bool
	StatusCode int
	Fail       int                 `json:"failure"`
	Results    []map[string]string `json:"results,omitempty"`
}

// Send send a single request to fcm
func (c *FCMClient) Send() (*FcmResponseStatus, error) {
	fcmRespStatus := new(FcmResponseStatus)

	jsonByte, err := json.Marshal(c.Message)
	if err != nil {
		fcmRespStatus.Ok = false
		return fcmRespStatus, err
	}

	revel.ERROR.Println("Notification: ", c.Message)

	// add headers tomessage
	request, _ := http.NewRequest("POST", FCMServerURL, bytes.NewBuffer(jsonByte))
	request.Header.Set("Authorization", "key="+getFCMServerKey())
	request.Header.Set("Content-Type", "application/json")

	// Create httpClient
	client := &http.Client{}

	// Send Request
	response, err := client.Do(request)
	if err != nil {
		fcmRespStatus.Ok = false
		return fcmRespStatus, err
	}

	revel.ERROR.Println("Response: ", response)

	// Close when finish function
	defer response.Body.Close()

	if _, err := ioutil.ReadAll(response.Body); err != nil {
		fcmRespStatus.Ok = false
		return fcmRespStatus, err
	}

	fcmRespStatus.Ok = true

	return fcmRespStatus, nil
}

// NewFCMClient ...
func NewFCMClient(to string, priority string, notification FCMNotification, ids []string, data interface{}) *FCMClient {
	var client = new(FCMClient)

	client.Message = FCMMessage{
		To:              to,
		RegistrationIds: ids,
		Priority:        priority,
		Data:            data,
		Notification:    notification,
		TimeToLive:      getTTL(),
	}

	return client
}

// NewFCMMessageTo ...
func (c *FCMClient) NewFCMMessageTo(to string, body interface{}) {
	c.Message.To = to
	c.Message.Data = body
}

// AppengRegIDs append the given IDS to the message
func (c *FCMClient) AppengRegIDs(ids []string) {
	c.Message.RegistrationIds = ids
}

// NewFcmRegIdsMsg gets a list of devices with data payload
func (c *FCMClient) NewFcmRegIdsMsg(list []string, body interface{}) *FCMClient {
	c.NewDevicesList(list)
	c.Message.Data = body

	return c
}

// NewDevicesList init the devices list
func (c *FCMClient) NewDevicesList(list []string) *FCMClient {
	c.Message.RegistrationIds = make([]string, len(list))
	copy(c.Message.RegistrationIds, list)

	return c

}

// AppendDevices adds more devices/tokens to the Fcm request
func (c *FCMClient) AppendDevices(list []string) *FCMClient {
	c.Message.RegistrationIds = append(c.Message.RegistrationIds, list...)
	return c
}

// NewSimpleFcmClient init and create fcm client
func NewSimpleFcmClient() *FCMClient {
	fcmc := new(FCMClient)
	return fcmc
}

// SetMessageData sets the data to the fcm cliente
func (c *FCMClient) SetMessageData(data interface{}) {
	c.Message.Data = data
}

// SetNotficationMessage appends the notification message
func (c *FCMClient) SetNotficationMessage(notification FCMNotification) *FCMClient {
	c.Message.Notification = notification
	return c
}

// getFCMServerKey returns fcm server key depend to env
func getFCMServerKey() string {
	if revel.RunMode != "prod" {
		return DevFCMServerKey
	}

	return FCMServerKey
}

// getTLL returns ttl to notifications depend to env
func getTTL() int {
	if revel.DevMode {
		return DevMaxTTL
	}

	return MAXTTL
}
