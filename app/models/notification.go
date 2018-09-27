package models

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"strconv"
	"strings"
	"time"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// Notification model
type Notification struct {
	mgomap.DocumentBase `json:"-" bson:",inline"`

	To          bson.ObjectId `bson:"to"`
	From        interface{}   `bson:"from"`
	Resource    interface{}   `bson:"resource"`
	ImageURL    string        `bson:"image_url"`
	Type        string        `bson:"type"`
	Action      string        `bson:"action"`
	Message     string        `bson:"message"`
	FullMessage string        `bson:"full_message"`
	Screen      string        `bson:"screen"`
	Title       string        `bson:"title"`
	Device      string        `bson:"device"`
	IDs         []string      `bson:"ids"`
	Attachment  Attachment    `bson:"attachment"`

	// Not saved fields
	ExtraData map[string]interface{} `json:"extra_data" bson:"-"`
}

// GetDocumentName ...
func (m *Notification) GetDocumentName() string {
	return "notifications"
}

// SendPush ...
func (m Notification) SendPush(to string, ids []string, badge int) (*core.FcmResponseStatus, error) {
	var from = m.ExtraData["from"]
	var resource = m.ExtraData["resource"]
	var notify core.FCMNotification

	if m.Device == "IOS" {
		notify.Title = m.Title
		notify.Body = m.FullMessage
		notify.Badge = strconv.Itoa(badge)
		notify.Sound = "default"
	}

	data := &struct {
		Type      string      `json:"type"`
		Id        string      `json:"id"`
		Image     string      `json:"image"`
		Message   string      `json:"message"`
		Action    string      `json:"action"`
		From      interface{} `json:"actor"`
		Resource  interface{} `json:"resource"`
		CreatedAt int64       `json:"created_at"`
	}{m.Type, m.GetID().Hex(), "", m.Message, m.Action, from, resource, m.CreatedAt.Unix()}

	if m.Attachment.PATH != "" {
		data.Image = m.Attachment.GetURL()
	}

	if resp, err := core.NewFCMClient(
		to,
		core.PriorityHigh,
		notify,
		ids,
		data,
	).Send(); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

func (m *Notification) SendToTopic(topic, id string, regIDs []string) {

	data := struct {
		Type     string      `json:"type"`
		Id       string      `json:"id"`
		Image    string      `json:"image"`
		Message  string      `json:"message"`
		Action   string      `json:"action"`
		From     interface{} `json:"actor"`
		Resource interface{} `json:"resource"`
	}{"notification", "id_notifica", "image_url", "Come and join us", "play", "from_id", id}

	c := core.FCMClient{}
	c.AppengRegIDs(regIDs)
	c.NewFCMMessageTo(topic, data)

	status, err := c.Send()

	if err == nil {
		revel.ERROR.Print(status)
	} else {
		revel.ERROR.Print(err)
	}

}

func (m *Notification) SendRegsIds(regIDs []string, device string) {
	fcm := core.NewSimpleFcmClient()
	var from = m.ExtraData["from"]
	var resource = m.Resource
	var notify core.FCMNotification

	// Create non-pointer struct to store as new Document in DB
	var notifica2 = Notification{
		Action:      m.Action,
		Device:      device,
		ExtraData:   m.ExtraData,
		From:        m.From,
		FullMessage: m.FullMessage,
		IDs:         regIDs,
		Message:     m.Message,
		Resource:    m.Resource,
		Screen:      m.Screen,
		Title:       m.Title,
		To:          m.To,
		Type:        m.Type,
		Attachment:  m.Attachment,
	}

	if device == "IOS" {
		notify.Title = m.Title
		notify.Body = m.FullMessage
		notify.Badge = "1"
		notify.Sound = "default"
		fcm.SetNotficationMessage(notify)
	}

	var data = &struct {
		Type      string      `json:"type"`
		Id        string      `json:"id"`
		Image     string      `json:"image"`
		Message   string      `json:"message"`
		Action    string      `json:"action"`
		From      interface{} `json:"actor"`
		Resource  interface{} `json:"resource"`
		CreatedAt interface{} `json:"created_at"`
	}{m.Type, m.Resource.(string), "", m.Message, m.Action, from, resource, time.Now().Unix()}

	if m.Attachment.PATH != "" {
		data.Image = m.Attachment.GetURL()
	}

	// Save Notification with the given registration IDs
	if Notifica, ok := app.Mapper.GetModel(&Notification{}); ok {
		if err := Notifica.Create(&notifica2); err != nil {
			revel.ERROR.Printf("ERROR CREATE Push Notification for Mission %s --- %s", m.Resource, err.Error())
		}
	}

	fcm.NewFcmRegIdsMsg(regIDs[0:1], data)
	if len(regIDs) > 1 {
		fcm.AppendDevices(regIDs[1:(len(regIDs) - 1)])
	}
	fcm.Send()
}

// Internationalize find translation for fields [FullMessage, Message, Title]
func (m *Notification) Internationalize(controller *revel.Controller) {
	var message = strings.Split(m.Message, "|")
	var args = make([]interface{}, len(message[1:]))
	for k, v := range message[1:] {
		args[k] = v
	}

	m.Message = strings.Replace(controller.Message(message[0], args...), "???", "", 2)

	message = strings.Split(m.Title, "|")
	for k, v := range message[1:] {
		args[k] = v
	}

	m.Title = strings.Replace(controller.Message(message[0], args...), "???", "", 2)
}
