package v2

import (
	"log"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// SearchController ...
type SearchController struct {
	BaseController
}

// TimeFormat returns the server timestampt for the request format
func (c SearchController) TimeFormat(format string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	t := time.Now()
	switch format {
	case "unix":
		timeStr := strconv.Itoa(int(int32(t.Unix())))
		return c.SuccessResponse(timeStr, "success", 0, nil)
	case "rfc3339":
		return c.SuccessResponse(t.Format(time.RFC3339), "success", 0, nil)
	case "latin":
		return c.SuccessResponse(t.Format("02-Jan-2006 15:04:05"), "success", 0, nil)
	case "europe":
		return c.SuccessResponse(t.Format("Jan-02-2006 15:04:05"), "success", 0, nil)
	default:
		return c.ErrorResponse("error", "error", 0)
	}
}

// GetPeople ...
func (c SearchController) GetPeople() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	people := &struct {
		People []string `json:"people"`
	}{}

	log.Print("Aqui 0")

	if err := c.PermitParams(people, false, "people"); err != nil {
		return c.ErrorResponse(err, err.Error(), 200)
	}
	log.Print("Aqui 1")
	// if User, ok := app.Mapper.GetModel(&c.CurrentUser); ok {
	var users = []models.User{}
	var match = bson.M{"device.number": bson.M{"$in": people.People}}
	var project = bson.M{
		"extra_parameters.are_friends": bson.M{"$in": []string{c.CurrentUser.GetID().Hex(), "$friends"}},
		"email":         1,
		"attachment":    1,
		"personal_data": 1,
		"device":        1,
	}
	var pipe = mgomap.Aggregate{}.Match(match).Add(bson.M{"$project": project})

	log.Print("Aqui 2")
	if err := models.AggregateQuery(&models.User{}, pipe, &users); err != nil {
		return c.ErrorResponse(err, err.Error(), 200)
	}

	log.Print("Aqui 3")

	return c.SuccessResponse(users, "People", 200, serializers.SearchPeopleSerializer{})
}

// Search ...
func (c SearchController) Search(account string, saas bool) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(account) < 1 || len(account) > 25 {
		return c.ErrorResponse(nil, "Account length not permited", 400)
	}

	// if User, ok := app.Mapper.GetModel(&c.CurrentUser); ok {
	var maxRegisters = core.DefaultLimitDocuments
	switch {
	case len(account) <= 3:
		maxRegisters = 15
	case len(account) > 3 && len(account) <= 5:
		maxRegisters = 25
	case len(account) > 5:
		maxRegisters = 50
	}

	var users []models.User
	var query []bson.M

	if saas {
		if err := models.Query(
			c.CurrentUser.GetDocumentName(),
			bson.M{"saas.name": bson.M{"$regex": account}}).Select(c.CurrentUser.GetPublicFields()).Exec(&users); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(users, "success", 200, serializers.UserSerializer{})
	}

	if _, err := strconv.Atoi(account); err != nil {
		query = append(query, bson.M{"email": account})
		query = append(query, bson.M{"user_name": bson.M{"$regex": account, "$options": "i"}})
		query = append(query, bson.M{"personal_data.first_name": bson.M{"$regex": account, "$options": "i"}})
		query = append(query, bson.M{"personal_data.last_name": bson.M{"$regex": account, "$options": "i"}})
	} else {
		query = append(query, bson.M{"device.number": bson.M{"$regex": account}})
	}

	var pipe = mgomap.Aggregate{}
	pipe = pipe.Match(bson.M{"$or": query}).Select(c.CurrentUser.GetPublicFields())
	pipe = pipe.Limit(maxRegisters)

	if err := models.AggregateQuery(models.User{}, pipe, &users); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(users, "Success", 200, serializers.UserSerializer{})
}

// InvitePeople ...
func (c SearchController) InvitePeople() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	people := &struct {
		People []string `json:"people"`
	}{}

	if err := c.PermitParams(people, false, "people"); err != nil {
		return c.ErrorResponse(err, err.Error(), 200)
	}

	for _, v := range people.People {
		// Enviamos el codigo SMS.
		if err := SendSMS(v, c.Message("sms.invite", c.CurrentUser.Email, "https://b8jt3.app.goo.gl/RtQw")); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	}

	return c.SuccessResponse("OK", "Sent", 200, nil)
}
