package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// SpyNewsController controller
type SpyNewsController struct {
	BaseController
}

// Index [GET, /spy_news]
// Returns all news for current user.
func (c SpyNewsController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var spyNews []models.SpyNews

	News, ok := app.Mapper.GetModel(&models.SpyNews{})
	if !ok {
		return c.ServerErrorResponse()
	}

	var pipe = mgomap.Aggregate{}.Match(
		bson.M{"$or": []bson.M{
			{"languages": bson.M{"$in": []string{c.Request.Locale}}},
			{"countries": bson.M{"$in": []string{c.Request.Header.Get("Country")}}},
		}},
	).Add(bson.M{"$addFields": bson.M{
		"extra.current_interaction": bson.M{"$filter": bson.M{"input": "$interactions", "as": "react",
			"cond": bson.M{"$eq": []interface{}{"$$react.user", c.CurrentUser.GetID()}}}}}},
	)

	if err := News.Pipe(pipe, &spyNews); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(spyNews, "Success", 200, serializers.SpyNewsSerializer{Lang: c.Request.Locale})
}

// Create [POST, /spy_news]
// Create a new SpyNews
func (c SpyNewsController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var spyNew models.SpyNews

	if err := c.PermitParams(&spyNew, true, "titles", "descriptions", "langs", "countries"); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if User, ok := app.Mapper.GetModel(&spyNew); ok {
		if err := User.Create(&spyNew); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(spyNew, "Success", 200, serializers.SpyNewsSerializer{Lang: c.Request.Locale})
	}

	return c.ServerErrorResponse()
}

// Comment create a comment for that spynews [POST, /spy_news/:id/comments]
func (c SpyNewsController) Comment(id, message string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var spynews models.SpyNews
	var comment models.Comment

	News, ok := app.Mapper.GetModel(&spynews)
	if !ok {
		return c.ServerErrorResponse()
	}

	if err := News.Find(id).Exec(&spynews); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if Comment, ok := app.Mapper.GetModel(&comment); ok {
		comment.Owner = c.CurrentUser.GetID()
		comment.Message = message
		comment.Target = spynews.GetDocumentName()
		comment.TargetID = spynews.GetID().Hex()

		if err := Comment.Create(comment); err != nil {
		}
	}

	return c.SuccessResponse("ok", "success", 200, nil)
}

// React [POST, /spy_news/:id/reactions]
func (c SpyNewsController) React(id, react string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var spynews models.SpyNews

	News, ok := app.Mapper.GetModel(&spynews)
	if !ok {
		return c.ServerErrorResponse()
	}

	total, err := News.Query(bson.M{
		"_id":          bson.ObjectIdHex(id),
		"interactions": bson.M{"$elemMatch": bson.M{"user": c.CurrentUser.GetID()}}}).Count()
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if total <= 0 {
		var selector = bson.M{"_id": bson.ObjectIdHex(id)}
		var query = bson.M{"$push": bson.M{"interactions": bson.M{"user": c.CurrentUser.GetID(), "name": react}}}
		if err := News.UpdateQuery(selector, query, false); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("ok", "success", 200, nil)
	}

	var selector = bson.M{"_id": bson.ObjectIdHex(id)}
	var query = bson.M{"$set": bson.M{"interactions.$.name": react}}
	if err := News.UpdateQuery(selector, query, false); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "success", 200, nil)
}
