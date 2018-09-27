package v2

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"spyc_backend/app/auth"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/go-ozzo/ozzo-validation"

	"encoding/json"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/gographics/imagick.v2/imagick"
	"gopkg.in/mgo.v2/bson"
)

// BaseController ...
type BaseController struct {
	*revel.Controller

	CurrentUser models.User
}

// CheckCompressHeader ...
func (c BaseController) CheckCompressHeader() bool {
	return c.Request.Header.Get("Accept-Encoding") == "application/gzip"
}

// SuccessResponse ...
func (c BaseController) SuccessResponse(data interface{}, message string, code int, s serializers.Serializer) revel.Result {
	c.Response.Status = http.StatusOK

	if c.Request.Header.Get(core.Compression) == core.TrueString {
		return serializers.SuccessResponse(data, message, code, s).Compress()
	}

	return serializers.SuccessResponse(data, message, code, s)
}

// ErrorResponse returns 400 code and json serializer
func (c BaseController) ErrorResponse(errors interface{}, message string, code int) revel.Result {
	c.Response.Status = http.StatusBadRequest

	if c.Request.Header.Get(core.Compression) == core.TrueString {
		return serializers.ErrorResponse(errors, message, code).Compress()
	}

	return serializers.ErrorResponse(errors, message, code)
}

// UnauthorizedResponse returns 401 code and json serializer
func (c BaseController) UnauthorizedResponse() revel.Result {
	c.Response.Status = http.StatusUnauthorized

	if c.Request.Header.Get(core.Compression) == core.TrueString {
		return serializers.UnauthorizedResponse().Compress()
	}

	return serializers.UnauthorizedResponse()
}

// ForbiddenResponse returns 403 code and json serializer
func (c BaseController) ForbiddenResponse() revel.Result {
	c.Response.Status = http.StatusForbidden

	if c.Request.Header.Get(core.Compression) == core.TrueString {
		return serializers.ForbiddenResponse().Compress()
	}

	return serializers.ForbiddenResponse()
}

// ServerErrorResponse returns 500 code and json serializer
func (c BaseController) ServerErrorResponse() revel.Result {
	c.Response.Status = 500

	if c.Request.Header.Get(core.Compression) == core.TrueString {
		return serializers.ServerErrorResponse().Compress()
	}

	return serializers.ServerErrorResponse()
}

// FileResponse ...
func (c BaseController) FileResponse(file []byte, id string) revel.Result {
	var path = "/tmp/" + id + ".png"

	if err := ioutil.WriteFile(path, file, 0644); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	f, err := os.Open(path)
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.RenderFile(f, revel.Inline)
}

// GetCurrentUser returns user's session
func (c *BaseController) GetCurrentUser() bool {
	if user, err := auth.AuthenticateByToken(c.Request); err == nil {
		c.CurrentUser = user
		return true
	}

	return false
}

// NewNotification ...
func (c BaseController) NewNotification(action string, attachment models.Attachment, message string, screen string, title string, t string, to models.User, resource interface{}, save bool) error {
	notification := models.Notification{
		From:        c.CurrentUser.UserName,
		To:          to.GetID(),
		Action:      action,
		Type:        t,
		Message:     message,
		Screen:      screen,
		Title:       title,
		FullMessage: message,
		Resource:    resource,
		Device:      c.CurrentUser.Device.OS,
		Attachment:  attachment,
	}

	var query = bson.M{
		"$and": []bson.M{
			{"to": to.GetID()},
			{"status.name": core.StatusInit},
		},
	}

	total, err := models.NumberOfDocuments(&notification, query)
	if err != nil {
		return err
	}

	if save {
		if err := models.CreateDocument(&notification); err != nil {
			return err
		}
	}

	if resp, err := notification.SendPush(to.Device.MessagingToken, []string{}, total+1); err != nil {
		return err
	} else {
		if !resp.Ok {
			revel.ERROR.Print(resp.Results)
		}
	}

	return errors.New("Model not found.\n")
}

// NewCustomNotification ...
// t = invitation, notification
func (c BaseController) NewCustomNotification(action string, attachment models.Attachment, message string, screen string, title string, t string, fromDocName string, to models.User, resource interface{}) error {

	notification := models.Notification{
		From:        fromDocName,
		To:          to.GetID(),
		Resource:    resource,
		Action:      action,
		Type:        t,
		Message:     message,
		FullMessage: message,
		Screen:      screen,
		Title:       title,
		Attachment:  attachment,
	}

	var query = bson.M{
		"$and": []bson.M{
			{"to": to.GetID()},
			{"status.name": core.StatusInit},
		},
	}

	total, err := models.NumberOfDocuments(&notification, query)
	if err != nil {
		return err
	}

	if err := models.CreateDocument(&notification); err != nil {
		return err
	}

	if resp, err := notification.SendPush(to.Device.MessagingToken, []string{}, total+1); err != nil {
		return err
	} else {
		if !resp.Ok {
			revel.ERROR.Print(resp.Results)
		}
	}

	return errors.New("Model not found")
}

// PermitParams ...
func (c BaseController) PermitParams(out interface{}, validate bool, allowed ...string) error {
	var params map[string]interface{}

	if err := c.Params.BindJSON(&params); err != nil {
		return err
	}

	log := c.Log.New("params", params)
	log.Info("Request")

	for k := range params {
		if core.FindOnArray(allowed, k) < 0 {
			delete(params, k)
		}
	}

	if j, err := json.Marshal(params); err != nil {
		return err
	} else {
		if err := json.Unmarshal(j, out); err != nil {
			return err
		}
	}

	if validate {
		if v, ok := out.(validation.Validatable); ok {
			return v.Validate()
		} else {
			return errors.New("Validatable object expected")
		}
	}

	return nil
}

// SendRegIds send the notification
func (c BaseController) SendRegIds(tokenList []string, device string, notification *models.Notification) {

	if len(tokenList) > 0 {
		var sliceTokens = make([]string, core.MaxPushRegIds)
		if len(tokenList) > core.MaxPushRegIds {
			sliceTokens = tokenList[0:core.MaxPushRegIds]
		} else {
			sliceTokens = tokenList[0:len(tokenList)]
		}
		if _, ok := notification.Resource.(string); ok {
			notification.SendRegsIds(sliceTokens, device)
			if len(tokenList) > core.MaxPushRegIds {
				c.SendRegIds(tokenList[(core.MaxPushRegIds+1):len(tokenList)], device, notification)
			}
		}
	}
}

// NotifyMissionComplete ...
func (c BaseController) NotifyMissionComplete(mission models.Mission) {
	var deviceList []struct {
		ID    string   `bson:"_id"`
		Users []string `bson:"users"`
	}

	var match = bson.M{"mission_id": mission.GetID(), "status.name": bson.M{"$nin": []string{core.StatusCompleted, core.StatusCompletedNotWinner, core.StatusCompleted}}}
	var unwind = bson.M{"$unwind": "$users"}
	var group = bson.M{"$group": bson.M{
		"_id":   "$users.device.os",
		"users": bson.M{"$push": "$users.device.messaging_token"},
	}}
	var pipe = mgomap.Aggregate{}.Match(match).LookUp("users", "user_id", "_id", "users").Add(unwind).Add(group)
	if err := models.AggregateQuery(&models.Mission{}, pipe, &deviceList); err != nil {
		revel.ERROR.Printf("ERROR FIND Users for mission: %s --- %s", mission.GetID().Hex(), err.Error())
		return
	}

	var notification = models.Notification{
		From:     mission.GetDocumentName(),
		To:       mission.GetID(),
		Resource: mission.GetID().Hex(),
		Action:   "play",
		Type:     core.CashHuntFinished,
		//Message:     c.Message("mission.complete", mission.Title),
		//FullMessage: c.Message("mission.complete", mission.Title),
		Message:     c.Message("mission.complete", serializers.InternationalizeSerializer{}.Get(mission.Title, c.Request.Locale)),
		FullMessage: c.Message("mission.complete", serializers.InternationalizeSerializer{}.Get(mission.Description, c.Request.Locale)),
		Screen:      "cash_hunt",
		Title:       c.Message("mission.completeTitle"),
		Device:      "",
		Attachment:  mission.Attachment,
	}

	if len(deviceList) == 0 {
		return
	}

	for i := 0; i < len(deviceList); i++ {
		if deviceList[i].ID == "Android" && len(deviceList[i].Users) > 0 {
			c.SendRegIds(deviceList[i].Users, "Android", &notification)
		}
		if deviceList[i].ID == "IOS" && len(deviceList[i].Users) > 0 {
			c.SendRegIds(deviceList[i].Users, "IOS", &notification)
		}
	}
}

// ResizeImage resizes a multipart to the given size using imagemagick
func (c BaseController) ResizeImage(size uint, multipart string, blur float64) []byte {
	mImage, err := c.Params.Files[multipart][0].Open()
	if err != nil {
		revel.ERROR.Printf("ERROR OPEN Multipart")
		return []byte{}
	}
	defer mImage.Close()
	buf := new(bytes.Buffer)

	if core.GetFileType(c.Params.Files[multipart][0].Filename) != "image" {
		revel.ERROR.Printf("ERROR UNSUPORTED FileType")
		return []byte{}
	}

	// Copy image into buffer
	if _, err := io.Copy(buf, mImage); err != nil {
		revel.ERROR.Printf("ERROR READ Buffer --- %s", err.Error())
		return []byte{}
	}

	byteArray := buf.Bytes()
	imagick.Initialize()

	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	// read image with ImageMagick
	err = mw.ReadImageBlob(byteArray)
	if err != nil {
		revel.ERROR.Printf("ERROR READ ImageMagick --- %s", err.Error())
		return []byte{}
	}

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(size, size, imagick.FILTER_LANCZOS, blur)
	if err != nil {
		revel.ERROR.Printf("ERROR RESIZE ImageMagick --- %s", err.Error())
		return []byte{}
	}

	// Set the compression quality to 95 (high quality = low compression)
	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		revel.ERROR.Printf("ERROR COMPRESSING ImageMagick --- %s", err.Error())
		return []byte{}
	}
	return mw.GetImageBlob()
}
