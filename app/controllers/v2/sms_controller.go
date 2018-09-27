package v2

import (
	"errors"
	"net/http"
	"net/url"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app/core"

	"github.com/revel/revel"
)

// SMSController ...
type SMSController struct {
	BaseController
}

// RequestNew ...
func (c *SMSController) RequestNew(phone string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if c.CurrentUser.Device.Status.Name != core.StatusInit && c.CurrentUser.Device.Status.Name != core.StatusPendingConfirmation {
		return c.ErrorResponse(nil, "The sms has been sent", 200)
	}

	if len(phone) >= 0 {
		c.CurrentUser.Device.Number = phone
		c.CurrentUser.Device.GenerateCode()
	} else {
		c.CurrentUser.Device.GenerateCode()
	}

	c.CurrentUser.Status.Name = core.StatusSMSSent
	c.CurrentUser.Status.Code = core.AccountStatus[c.CurrentUser.Status.Name]

	var fields = bson.M{
		"device.number": c.CurrentUser.Device.Number,
		"device.code":   c.CurrentUser.Device.Code,
		"device.status": c.CurrentUser.Device.Status,
		"status":        c.CurrentUser.Status,
	}

	if err := models.SetDocument(c.CurrentUser.GetDocumentName(), c.CurrentUser.GetID().Hex(), fields); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	// Enviamos el codigo SMS.
	m := c.Message("sms.verification", c.CurrentUser.Device.Code)
	if err := SendSMS(c.CurrentUser.Device.Number, m); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if err := core.SendEmail([]string{c.CurrentUser.Email}, "Codigo de confirmacion spychatter", m); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("Sent", "Sent SMS", 200, nil)
}

// VerifySMS codig por SMS
func (c *SMSController) VerifySMS(code string) revel.Result {
	if !c.GetCurrentUser() {
		c.ForbiddenResponse()
	}

	if c.CurrentUser.Device.Code == code &&
		c.CurrentUser.Device.Status.Name == core.StatusPendingConfirmation {

		// if User, ok := app.Mapper.GetModel(&c.CurrentUser); ok {
		c.CurrentUser.Device.Status.Name = core.StatusActive
		c.CurrentUser.Device.Status.Code = core.AccountStatus[c.CurrentUser.Device.Status.Name]
		c.CurrentUser.Device.RandomCode()
		c.CurrentUser.SetState(core.StatusSMSConfirmed)
		c.CurrentUser.Saas.Init(c.CurrentUser.Email, c.CurrentUser.Device.Number)

		var fields = bson.M{
			"device.number": c.CurrentUser.Device.Number,
			"device.code":   c.CurrentUser.Device.Code,
			"device.status": c.CurrentUser.Device.Status,
			"status":        c.CurrentUser.Status,
			"saas":          c.CurrentUser.Saas,
		}

		if err := models.SetDocument(c.CurrentUser.GetDocumentName(), c.CurrentUser.GetID().Hex(), fields); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(c.CurrentUser, "Confirmed!!", 200, serializers.UserSerializer{})
	}

	return c.ErrorResponse(nil, "Invalid code", 400)
}

// SendSMS is a function for send SMS message to number.
// This message contains a code for user account activation.
// Number fromat: (+[0-9,2])[0-9,10]
// Code format: ([0-9, 6])
func SendSMS(number string, message string) error {
	// Twilio credentials
	var accountSID = revel.Config.StringDefault("spy.secrets.twilio.sid", "")
	var authToken = revel.Config.StringDefault("spy.secrets.twilio.token", "")
	var fromNumber = revel.Config.StringDefault("spy.secrets.twilio.number", "")
	var twilioURL = "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"

	revel.INFO.Printf("SENDING %s TO %s", message, number)

	values := url.Values{}
	values.Set("To", number)
	values.Set("From", fromNumber)
	values.Set("Body", message)

	rb := *strings.NewReader(values.Encode())

	client := &http.Client{}

	req, err := http.NewRequest("POST", twilioURL, &rb)
	if err != nil {
		return err
	}

	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	revel.INFO.Println(resp.Status[0:3])
	if resp.Status[0:3] != "201" {
		return errors.New("Error to try send SMS to " + number)
	}

	return nil
}
