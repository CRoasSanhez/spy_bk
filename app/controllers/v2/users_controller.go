package v2

import (
	"log"
	"regexp"
	"spyc_backend/app"
	"spyc_backend/app/auth"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// UsersController controller for Users
type UsersController struct {
	BaseController
}

// Index [GET, /profile]
// returns current user information
func (c UsersController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	return c.SuccessResponse(c.CurrentUser, "success", 200, serializers.UserSerializer{})
}

// Show [GET, /users/:id]
// function to get user profile
func (c UsersController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if id == c.CurrentUser.GetID().Hex() {
		return c.SuccessResponse(c.CurrentUser, "success", 200, serializers.UserSerializer{})
	}

	var user models.User
	if err := user.PublicProfile(id); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(user, "Success", 200, serializers.UserSerializer{})
}

// Create [POST, /users]
// Generate a new spychatter account
func (c UsersController) Create(body string) revel.Result {
	var user models.User

	if err := c.PermitParams(&user, true, "email", "password", "geolocation", "personal_data", "device"); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	user.Email = strings.ToLower(user.Email)
	user.Device.Init()
	user.SetState(core.StatusCreated)
	user.GeneratePassword()
	user.Geolocation.Init()

	if err := models.CreateDocument(&user); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	models.CreateStats(user.GetID())
	models.SetSections(user.GetID())

	token, err := auth.GenerateToken(user.GetID().Hex(), core.ActionAuth)
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(user, "Created!!", 200, serializers.UserSerializer{Token: token})
}

// Update [PATCH, /profile]
// Function to update user info, only [user_name, "personal_data", "device"]
func (c UsersController) Update() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var user models.User

	if err := c.PermitParams(&user, false, "email", "user_name", "personal_data", "device"); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if err := user.PersonalData.Validate(); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	core.ClearFields(&user.Device, "Number", "Code", "Status")
	core.MergeStructs(&c.CurrentUser.Device, &user.Device)
	core.MergeStructs(&c.CurrentUser.PersonalData, &user.PersonalData)

	if user.Email != "" && user.Email != c.CurrentUser.Email {
		if total, err := models.NumberOfDocuments(&c.CurrentUser, bson.M{"email": user.Email}); total > 0 || err != nil {
			return c.ErrorResponse(nil, "Duplicate key", 1011)
		}

		c.CurrentUser.Email = user.Email
	}

	var fields = bson.M{
		"personal_data": c.CurrentUser.PersonalData,
		"device":        c.CurrentUser.Device,
		"email":         c.CurrentUser.Email,
	}

	if len(c.CurrentUser.UserName) <= 0 {
		fields["user_name"] = strings.ToLower(user.UserName)
	}

	if err := models.SetDocument(c.CurrentUser.GetDocumentName(), c.CurrentUser.GetID().Hex(), fields); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(c.CurrentUser, "Updated!!", 200, serializers.UserSerializer{})
}

// Delete [DELETE, /profile]
func (c UsersController) Delete() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var fields = bson.M{"status.name": core.StatusSuspended, "status.code": core.AccountStatus[core.StatusSuspended]}

	if err := models.SetDocument(c.CurrentUser.GetDocumentName(), c.CurrentUser.GetID().Hex(), fields); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "success", 200, nil)
}

// ProfilePicture upload an image to S3 and set as profile picture
func (c UsersController) ProfilePicture() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(c.Params.Files["picture"]) <= 0 {
		return c.ErrorResponse(nil, "File not found", 400)
	}

	if c.CurrentUser.Attachment.PATH != "" {
		c.CurrentUser.Attachment.Remove()
	}

	if err := c.CurrentUser.Attachment.Init(models.AsDocumentBase(&c.CurrentUser), c.Params.Files["picture"][0]); err != nil {
		return c.ErrorResponse(err, err.Error(), 402)
	}

	if err := c.CurrentUser.Attachment.Upload(); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	var fields = bson.M{"attachment": c.CurrentUser.Attachment}
	if err := models.SetDocument(c.CurrentUser.GetDocumentName(), c.CurrentUser.GetID().Hex(), fields); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(c.CurrentUser.Attachment, "Updated!!", 200, serializers.AttachmentSerializer{})
}

// ChangeUsername update a username of current user
func (c UsersController) ChangeUsername(username string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	username = strings.ToLower(username)

	var countQuery = bson.M{"user_name": username}
	if total, err := models.NumberOfDocuments(&c.CurrentUser, countQuery); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	} else if total > 0 {
		return c.ErrorResponse(nil, "Username is not available", 400)
	}

	var selector = bson.M{"_id": c.CurrentUser.GetID()}
	var query = bson.M{"$set": bson.M{"user_name": username}}

	if err := models.UpdateByQuery(&c.CurrentUser, selector, query, false); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "updated!!", 200, nil)
}

// Language update a user default language
func (c UsersController) Language(lang string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(lang) <= 0 {
		return c.ErrorResponse(nil, "lang error", 400)
	}

	var selector = bson.M{"_id": c.CurrentUser.GetID()}
	var query = bson.M{"$set": bson.M{"device.language": lang}}
	if err := models.UpdateByQuery(&c.CurrentUser, selector, query, false); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "updated!!", 200, nil)
}

// Coords function updates a user location
func (c UsersController) Coords(lat, lng float64) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// TODO: Check if coordinates are valid.
	var geo = models.Geo{
		Type:        "Point",
		Coordinates: []float64{lng, lat},
	}

	var fields = bson.M{"geolocation": geo}
	if err := models.SetDocument(c.CurrentUser.GetDocumentName(), c.CurrentUser.GetID().Hex(), fields); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "updated!!", 200, nil)
}

// MessagingToken update token for push notifications
func (c UsersController) MessagingToken(token string) revel.Result {
	if !c.GetCurrentUser() || len(token) <= 0 {
		return c.ForbiddenResponse()
	}

	var fields = bson.M{"device.messaging_token": token}
	if err := models.SetDocument(c.CurrentUser.GetDocumentName(), c.CurrentUser.GetID().Hex(), fields); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("OK", "Updated!!", 200, nil)
}

// VerifyUserName find user by username, return 200 if user exists.
func (c UsersController) VerifyUserName(account string) revel.Result {
	c.GetCurrentUser()

	if len(account) < 3 {
		return c.ErrorResponse(nil, "The user_name "+account+" is not available", 400)
	}

	account = strings.ToLower(account)

	if c.CurrentUser.UserName == account {
		return c.SuccessResponse("Equal", "The username is the current", 200, nil)
	}

	// TODO: Check this comparation
	if total, err := models.NumberOfDocuments(&c.CurrentUser, bson.M{"user_name": account}); total <= 0 || err != nil {
		return c.SuccessResponse("OK", "Available", 200, nil)
	}

	return c.ErrorResponse(nil, "The user_name "+account+" is not available", 400)
}

// FriendRequest ...
func (c UsersController) FriendRequest(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// User collection instance
	// if User, ok := app.Mapper.GetModel(&c.CurrentUser); ok {
	// friend
	var user2 models.User

	if err := models.GetDocument(id, &user2); err != nil {
		log.Println("ERROR: ", err.Error())
		return c.ErrorResponse(err, err.Error(), 400)
	}

	// ---------------------------------------------------------------------------
	// Obtain number of invitations
	// Find if exists an invite
	var invite models.Invitation

	var inviteQuery = bson.M{"$and": []bson.M{
		bson.M{"sender.$id": c.CurrentUser.GetID()},
		bson.M{"receiver.$id": user2.GetID()},
		bson.M{"type": "friend"},
	}}

	if total, err := models.NumberOfDocuments(&invite, inviteQuery); err != nil {
		log.Println("ERROR: ", err.Error())
		return c.ErrorResponse(err, err.Error(), 400)
	} else if total > 0 {
		log.Println("ERROR: ", err.Error())
		return c.ErrorResponse(nil, "Invite has been sent", 400)
	}

	// ---------------------------------------------------------------------------
	// ---------------------------------------------------------------------------
	// Create an invite
	var sender = &mgo.DBRef{
		Id:         c.CurrentUser.GetID(),
		Collection: c.CurrentUser.GetDocumentName(),
		Database:   app.Mapper.DatabaseName,
	}

	var receiver = &mgo.DBRef{
		Id:         user2.GetID(),
		Collection: user2.GetDocumentName(),
		Database:   app.Mapper.DatabaseName,
	}

	invite = models.Invitation{
		Sender:       sender,
		Receiver:     receiver,
		Resource:     receiver,
		ResourceType: user2.GetDocumentName(),
		Type:         "friend",
	}

	if err := models.CreateDocument(&invite); err != nil {
		log.Println("ERROR: ", err.Error())
		return c.ErrorResponse(err, err.Error(), 400)
	}

	// ---------------------------------------------------------------------------
	// ---------------------------------------------------------------------------
	// TODO make push notifications
	// Notification
	if lang := core.FindOnArray(core.CurrentLocales, c.CurrentUser.Device.Language); lang >= 0 {
		c.Request.Locale = c.CurrentUser.Device.Language
	}

	var title = c.Message("notification.new_friend_request_title")
	var message = c.Message("notification.new_friend_request", c.CurrentUser.UserName)
	c.NewNotification("friend request", c.CurrentUser.Attachment, message, "invitations", title, core.FriendRequest, user2, nil, false)

	return c.SuccessResponse("OK", "Sent", 200, nil)
}

// Friends ...
func (c UsersController) Friends() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	return c.SuccessResponse(c.CurrentUser.Friends, "success", 200, serializers.FriendSerializer{})
}

// ResetPassword ...
func (c UsersController) ResetPassword(email string) revel.Result {
	var user models.User
	var field = "email"

	isDevice, _ := regexp.MatchString("\\+[0-9]+", email)

	if isDevice {
		field = "device.number"
	}

	if err := models.GetDocumentBy(bson.M{field: email}, &user); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if lang := core.FindOnArray(core.CurrentLocales, user.Device.Language); lang >= 0 {
		c.Request.Locale = user.Device.Language
	}

	var message string

	token, err := models.NewToken(core.ActionAuth, user.GetDocumentName(), user.GetID(), time.Now().Add(time.Hour*168))
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	message = c.Message("email.reset_password", user.FullName(), core.GetDashboardPath()+"reset_password?token="+token)

	if isDevice {
		message = c.Message("email.reset_password_sms", user.FullName(), core.GetDashboardPath()+"reset_password?token="+token)
		SendSMS(user.Device.Number, message)
	} else {
		if err := core.SendEmail([]string{user.Email}, "Reset password", message); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	}

	return c.SuccessResponse("OK", "Email sent", 200, nil)
}
