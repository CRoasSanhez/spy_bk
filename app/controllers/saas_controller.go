package controllers

import (
	"spyc_backend/app"
	"spyc_backend/app/controllers/v2"
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"gopkg.in/mgo.v2/bson"

	"errors"

	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/revel/revel"
)

// SaasController controller
type SaasController struct {
	v2.BaseController
}

// SaaSStanza stanza for xmmp message
type SaaSStanza struct {
	UserName string `json:"username"`
	Peer     string `json:"peer"`
	Type     string `json:"type"`
	Nick     string `json:"nick"`
	XML      string `json:"xml"`
}

// Body stanza body
type Body struct {
	XMLName xml.Name `xml:"body"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// Message stanza message
type Message struct {
	XMLName xml.Name `xml:"message"`
	From    string   `xml:"from,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"`
	Id      string   `xml:"id,attr"`
	Body    []Body   `xml:"body"`
}

// Auth  for staless mode integration
func (c SaasController) Auth() revel.Result {
	username := c.Params.Get("username")
	password := c.Params.Get("password")

	if len(username) <= 0 || len(password) <= 0 {
		return c.UnauthorizedResponse()
	}

	var user models.User
	var query = bson.M{
		"$and": []bson.M{
			{"saas.name": username},
			{"saas.token": password},
		},
	}

	if err := models.GetDocumentBy(query, &user); err != nil {
		return c.UnauthorizedResponse()
	}

	return c.SuccessResponse("OK", "Success", 200, nil)
}

// User ...
func (c SaasController) User() revel.Result {
	username := c.Params.Get("username")

	if len(username) <= 0 {
		return c.UnauthorizedResponse()
	}

	var user models.User
	var query = bson.M{"saas.name": username}

	if err := models.GetDocumentBy(query, &user); err != nil {
		return c.UnauthorizedResponse()
	}

	return c.SuccessResponse("OK", "Success", 200, nil)
}

// Archive save a message on database
func (c SaasController) Archive() revel.Result {
	var stanza SaaSStanza
	if err := c.Params.BindJSON(&stanza); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	stanza.XML = strings.Replace(stanza.XML, "xml:lang", "lang", 0)

	var message Message
	if err := xml.Unmarshal([]byte(stanza.XML), &message); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	var from models.User
	var to models.User

	if err := models.GetDocumentBy(bson.M{"saas.name": stanza.UserName}, &from); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	// if err := User.FindBy("saas.name", stanza.Peer).Exec(&to); err != nil {
	if err := models.GetDocumentBy(bson.M{"saas.name": stanza.Peer}, &to); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	var body string
	var customBody map[string]interface{}

	for _, v := range message.Body {
		if v.Lang == "custom" {
			if err := json.Unmarshal([]byte(v.Value), &customBody); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {
			body = v.Value
		}
	}

	if len(body) <= 0 {
		body = "Has new messages"
	}

	var sendN = true

	switch customBody["elementType"] {
	case "IMAGE":
		body = c.Message("notification.image_message")
	case "DOCUMENT":
		body = c.Message("notification.document_message")
	case "LOCATION":
		body = c.Message("notification.location_message")
	case "CHALLENGE":
		body = c.Message("notification.challenge_message")
	case "GIF":
		body = c.Message("notification.gif_message")
	case "ANIMATION":
		body = c.Message("notification.animation_message")
	case "STICKER":
		body = c.Message("notification.sticker_message")
	case "AUDIO":
		body = c.Message("notification.audio_message")
	case "VIDEO":
		body = c.Message("notification.video_message")
	case "LINK":
		body = c.Message("notification.link_message")
	case "ANSWERCHALLENGE":
		sendN = false
	case "CONTACT":
		body = c.Message("notification.contact_message")
	case "CONFIRM":
		sendN = false
	case "DELETE":
		sendN = false
	case "STATUSCHALLENGE":
		body = c.Message("notification.status_challenge_message")
	case "SCREENSHOT":
		body = c.Message("notification.screenshot_message")
	case "SCHEDULE":
		sendN = false
	case "ONLINE":
		sendN = false
	default:
		sendN = true
	}

	var toJID = strings.Split(message.To, "@")[0]

	var schedule bool
	if flag, ok := customBody["schedule"].(string); ok {
		schedule = flag == "true"
	}

	if sendN && (toJID == stanza.Peer) && !schedule {
		var chatMessage = models.Message{
			Sender:     from.Saas.Name,
			Receiver:   to.Saas.Name,
			EjabberdID: message.Id,
			Message:    body,
			Stanza:     stanza.XML,
		}

		if err := models.CreateDocument(&chatMessage); err != nil {
			revel.ERROR.Print(err)
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if lang := core.FindOnArray(core.CurrentLocales, to.Device.Language); lang >= 0 {
			c.Request.Locale = to.Device.Language
		}

		var message = c.Message("notification.new_message_from_title", from.UserName)
		c.NewNotification("new message", body, "chat", message, core.Message, to, from)
	}

	return c.SuccessResponse("OK", "Success", 200, nil)
}

// GetRoster ...
// TODO(x):
func (c SaasController) GetRoster() revel.Result {
	username := c.Params.Get("username")

	if username == "" {
		return c.UnauthorizedResponse()
	}

	return c.SuccessResponse("Ok", "Success", 200, nil)
}

// NewRoster ...
// TODO(x):
func (c SaasController) NewRoster() revel.Result {
	return c.SuccessResponse("Ok", "Success", 200, nil)
}

// UpdateRoster ....
// TODO(x):
func (c SaasController) UpdateRoster() revel.Result {
	return c.SuccessResponse("Ok", "Success", 200, nil)
}

// NewNotification ...
// TODO(x):
func (c SaasController) NewNotification(action string, message string, screen string, title string, t string, to models.User, from models.User) error {
	notification := models.Notification{
		To:          to.GetID(),
		From:        from.UserName,
		Action:      action,
		Type:        t,
		Message:     message,
		Screen:      screen,
		Title:       title,
		FullMessage: message,
		Device:      to.Device.OS,
	}

	if Notification, ok := app.Mapper.GetModel(&notification); ok {
		var query = []bson.M{
			{"to": to.GetID()},
			{"status.name": core.StatusInit},
		}

		total, err := Notification.FindWithOperator("$and", query).Count()
		if err != nil {
			return err
		}

		if err := Notification.Create(&notification); err != nil {
			return err
		}

		if _, err := notification.SendPush(to.Device.MessagingToken, []string{}, total+1); err != nil {
			return err
		}
	}

	return errors.New("Model not found")
}
