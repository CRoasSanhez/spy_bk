package dashboard

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"spyc_backend/app"
	"spyc_backend/app/auth"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/gographics/imagick.v2/imagick"
	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// DBaseController controller
type DBaseController struct {
	*revel.Controller

	CurrentUser models.User
}

// SuccessResponse returns 200 code and json serializer
func (c DBaseController) SuccessResponse(data interface{}, message string, code int, s serializers.Serializer) revel.Result {
	c.Response.Status = http.StatusOK

	return c.RenderJSON(serializers.SuccessResponse(data, message, code, s))
}

// ErrorResponse returns 400 code and json serializer
func (c DBaseController) ErrorResponse(errors interface{}, message string, code int) revel.Result {
	c.Response.Status = http.StatusBadRequest

	return c.RenderJSON(serializers.ErrorResponse(errors, message, code))
}

// UnauthorizedResponse returns 401 code and json serializer
func (c DBaseController) UnauthorizedResponse() revel.Result {
	c.Response.Status = http.StatusUnauthorized

	return c.RenderJSON(serializers.UnauthorizedResponse())
}

// ForbiddenResponse returns 403 code and json serializer
func (c DBaseController) ForbiddenResponse() revel.Result {
	c.Response.Status = http.StatusForbidden

	return c.RenderJSON(serializers.ForbiddenResponse())
}

// ServerErrorResponse returns 500 code and json serializer
func (c DBaseController) ServerErrorResponse() revel.Result {
	c.Response.Status = 500

	return c.RenderJSON(serializers.ServerErrorResponse())
}

// FileResponse ...
func (c DBaseController) FileResponse(file []byte, id string) revel.Result {
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

// GetCurrentUser ...
func (c DBaseController) GetCurrentUser() bool {
	if c.Params.Get("format") == core.JSONFormat {
		// Login by token
		user, err := auth.AuthenticateByToken(c.Request)
		if err != nil {
			return false
		}

		c.CurrentUser = user

		return true
	}

	// Login by session cookie
	user, err := auth.AuthenticateBySession(c.Session[core.AccessToken])
	if err != nil || !user.CanManageSection(core.SectionDashboard) {
		return false
	}

	c.ViewArgs["user"] = user

	return true
}

// NewNotification ...
// t = invitation, notification
func (c DBaseController) NewNotification(action string, attachment models.Attachment, message string, screen string, title string, t string, fromDocName string, to models.User, resource interface{}) error {

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
		Device:      to.Device.OS,
		Attachment:  attachment,
	}

	Notification, _ := app.Mapper.GetModel(&notification)

	if _, ok := app.Mapper.GetModel(&notification); ok {
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

		if resp, err := notification.SendPush(to.Device.MessagingToken, []string{}, total+1); err != nil {
			return err
		} else {
			if !resp.Ok {
				revel.ERROR.Print(resp.Results)
			}
		}
	}

	return errors.New("Model not found")
}

// PermitParams for the given struct
func (c DBaseController) PermitParams(out interface{}, allowed ...string) error {
	var params map[string]interface{}

	if err := c.Params.BindJSON(&params); err != nil {
		return err
	}

	for k, _ := range params {
		if !FindOnArray(allowed, k) {
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

	return nil
}

// NotifyMissionComplete ...
func (c DBaseController) NotifyMissionComplete(mission models.Mission) {
	revel.INFO.Print("Notifying Mission complete")

	if MissionUser, ok := app.Mapper.GetModel(&models.MissionUser{}); ok {

		var err error
		var deviceList []struct {
			ID    string   `bson:"_id"`
			Users []string `bson:"users"`
		}

		c.Request.Locale = mission.Langs[0]

		var match = bson.M{"mission_id": mission.GetID(), "status.name": bson.M{"$nin": []string{core.StatusCompleted, core.StatusCompletedNotWinner, core.StatusCompleted}}}
		var unwind = bson.M{"$unwind": "$users"}
		var group = bson.M{"$group": bson.M{
			"_id":   "$users.device.os",
			"users": bson.M{"$push": "$users.device.messaging_token"},
		}}
		var pipe = mgomap.Aggregate{}.Match(match).LookUp("users", "user_id", "_id", "users").Add(unwind).Add(group)
		if err = MissionUser.Pipe(pipe, &deviceList); err != nil {
			revel.ERROR.Printf("ERROR FIND Users for mission: %s --- %s", mission.GetID().Hex(), err.Error())
			return
		}

		var notification = models.Notification{
			From:        mission.GetDocumentName(),
			To:          mission.GetID(),
			Resource:    mission.GetID().Hex(),
			Action:      "play",
			Type:        core.CashHuntFinished,
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

		// Send to android and ios respectlively
		for i := 0; i < len(deviceList); i++ {
			if deviceList[i].ID == "Android" && len(deviceList[i].Users) > 0 {
				SendRegIds(deviceList[i].Users, "Android", &notification)
			}
			if deviceList[i].ID == "IOS" && len(deviceList[i].Users) > 0 {
				SendRegIds(deviceList[i].Users, "IOS", &notification)
			}
		}
		return
	}
}

// SendRegIds send the notification
func (c DBaseController) SendRegIds(tokenList []string, device string, notification *models.Notification) {

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

// ResizeImage resizes a multipart to the given size using imagemagick
func (c DBaseController) ResizeImage(size uint, multipart string, blur float64) []byte {
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

func FindOnArray(array []string, key string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == key {
			return true
		}
	}

	return false
}
