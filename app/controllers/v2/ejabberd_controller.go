package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// EjabberdController ...
type EjabberdController struct {
	BaseController
}

// Index [GET] returns all groups of the current user
func (c EjabberdController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if Groups, ok := app.Mapper.GetModel(&models.EjabberdGroup{}); ok {
		var groups []models.EjabberdGroup
		var query mgomap.Query

		if len(c.Params.Get("jid")) > 0 {
			query = Groups.FindBy("jid", c.Params.Get("jid"))
		} else {
			query = Groups.FindBy("owner", c.CurrentUser.Saas.Name)
		}

		if err := query.Exec(&groups); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(groups, "Success", 200, serializers.EjabberdGroupSerializer{})
	}

	return c.ServerErrorResponse()
}

// Show []
func (c EjabberdController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var group models.EjabberdGroup
	if Ejabberd, ok := app.Mapper.GetModel(&group); ok {
		if err := Ejabberd.Find(id).Exec(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(group, "Success!!", 200, serializers.EjabberdGroupSerializer{})
	}

	return c.ErrorResponse(nil, "Invalid group", 403)
}

// Create function
func (c EjabberdController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var group models.EjabberdGroup
	if err := c.PermitParams(&group, true, "name", "jid", "members"); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	group.Owner = c.CurrentUser.Saas.Name

	if Ejabber, ok := app.Mapper.GetModel(&group); ok {
		if err := Ejabber.Create(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(group, "Created!!", 200, serializers.EjabberdGroupSerializer{})
	}

	return c.ServerErrorResponse()
}

func (c EjabberdController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var group models.EjabberdGroup

	if Ejabberd, ok := app.Mapper.GetModel(&group); ok {
		if err := Ejabberd.Find(id).Exec(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if len(c.Params.Get("name")) > 0 {
			group.Name = c.Params.Get("name")
		}

		if len(c.Params.Files["picture"]) > 0 {
			if group.Attachment.PATH != "" {
				group.Attachment.Remove()
			}

			if err := group.Attachment.Init(models.AsDocumentBase(&c.CurrentUser), c.Params.Files["picture"][0]); err != nil {
				return c.ErrorResponse(err, err.Error(), 402)
			}

			if err := group.Attachment.Upload(); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}
		}

		if err := Ejabberd.Update(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(group, "Updated!!", 200, serializers.EjabberdGroupSerializer{})
	}

	return c.ServerErrorResponse()
}

// AddMember ...
func (c EjabberdController) AddMember(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var group models.EjabberdGroup
	var query = []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"owner": c.CurrentUser.Saas.Name},
	}

	if Ejabberd, ok := app.Mapper.GetModel(&group); ok {
		if err := Ejabberd.FindWithOperator("$and", query).Exec(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		var newMember models.EjabberdUser
		newMember.JID = c.Params.Get("jid")
		newMember.Nick = c.Params.Get("nick")
		newMember.Admin = c.Params.Get("admin") == "true"

		if len(newMember.JID) <= 0 {
			return c.ErrorResponse(nil, "Invalid member", 400)
		}

		group.AddMember(newMember.JID, newMember.Nick, newMember.Admin)
		Ejabberd.Update(&group)

		return c.SuccessResponse(group, "Success", 200, serializers.EjabberdGroupSerializer{})
	}

	return c.ErrorResponse(nil, "Server error", 500)
}

func (c EjabberdController) UpdateMember(id, jid string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var group models.EjabberdGroup
	var query = []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"owner": c.CurrentUser.Saas.Name},
	}

	if Ejabberd, ok := app.Mapper.GetModel(&group); ok {
		if err := Ejabberd.FindWithOperator("$and", query).Exec(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		var member models.EjabberdUser
		if err := c.PermitParams(&member, true, "jid", "nick", "admin"); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		for _, v := range group.Members {
			if v.JID == jid {
				v.JID = member.JID
				v.Nick = member.Nick
				v.Admin = member.Admin
			}
		}

		if err := Ejabberd.Update(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(group, "Success", 200, serializers.EjabberdGroupSerializer{})
	}

	return c.ErrorResponse(nil, "Server error", 500)
}

// RemoveMember ...
func (c EjabberdController) RemoveMember(id, jid string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var group models.EjabberdGroup
	var query = []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"owner": c.CurrentUser.Saas.Name},
	}

	if Ejabberd, ok := app.Mapper.GetModel(&group); ok {
		if err := Ejabberd.FindWithOperator("$and", query).Exec(&group); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if len(jid) <= 0 {
			return c.ErrorResponse(nil, "Invalid member", 400)
		}

		group.RemoveMember(jid)
		Ejabberd.Update(&group)

		return c.SuccessResponse(group, "Success", 200, serializers.EjabberdGroupSerializer{})
	}

	return c.ErrorResponse(nil, "Server error", 500)
}

func (c EjabberdController) ExistsMember(members []models.EjabberdUser, jid string) interface{} {
	for _, v := range members {
		if v.JID == jid {
			return v
		}
	}

	return nil
}
