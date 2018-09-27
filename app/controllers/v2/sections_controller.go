package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// SectionsController controller
type SectionsController struct {
	BaseController
}

// Create add a section to current user
func (c SectionsController) Create(t int, name string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var s = models.CommonSection{
		Name:       name,
		AccessTime: t,
	}

	if Common, ok := app.Mapper.GetModel(&s); ok {
		if err := Common.Create(&s); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(s, core.StatusSuccess, 200, serializers.CommonSectionSerializer{})
	}

	return c.ServerErrorResponse()
}

// Replicate ...
func (c SectionsController) Replicate(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	if Common, ok := app.Mapper.GetModel(&models.CommonSection{}); ok {
		var section models.CommonSection

		if err := Common.Find(id).Exec(&section); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if err := section.Replicate(); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("ok", core.StatusSuccess, 200, nil)
	}

	return c.ServerErrorResponse()
}

// Info return general data of selected section
func (c SectionsController) Info(section string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	for _, v := range c.CurrentUser.Sections {
		if v.Name == section {
			switch section {
			case core.SectionCashHunt:
				if v.LastAccess.Add(v.AccessTime).Unix() < time.Now().Unix() {
					v.CanAccess = false
					v.TimeLess = 0

					if err := c.UpdateSection(c.CurrentUser.GetID(), v); err != nil {
						return c.ErrorResponse(err, err.Error(), 400)
					}

					return c.SuccessResponse(v, core.StatusSuccess, 200, serializers.SectionSerializer{})
				}

				v.TimeLess = v.LastAccess.Add(v.AccessTime).Sub(time.Now())
				return c.SuccessResponse(v, core.StatusSuccess, 200, serializers.SectionSerializer{})
			}
		}
	}

	return c.SuccessResponse(models.Section{Name: section, CanAccess: true, LastAccess: time.Now()}, core.StatusSuccess, 200, nil)
}

// NotifyAdvertisingSeen update section data
func (c SectionsController) NotifyAdvertisingSeen(section string) revel.Result {
	if !c.GetCurrentUser() {
		c.ForbiddenResponse()
	}

	for _, v := range c.CurrentUser.Sections {
		if v.Name == section {
			v.CanAccess = true
			v.LastAccess = time.Now()
			v.TimeLess = time.Now().Add(v.AccessTime).Sub(time.Now())

			if err := c.UpdateSection(c.CurrentUser.GetID(), v); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}

			return c.SuccessResponse(v, core.StatusSuccess, 200, serializers.SectionSerializer{})
		}
	}

	return c.ErrorResponse(nil, "section not found", 400)
}

// UpdateSection remove section with name [Name] and add section
func (c SectionsController) UpdateSection(id bson.ObjectId, s models.Section) error {
	if User, ok := app.Mapper.GetModel(&c.CurrentUser); ok {
		var selector = bson.M{"_id": id}
		var query = bson.M{"$pull": bson.M{"sections": bson.M{"name": s.Name}}}

		if err := User.UpdateQuery(selector, query, true); err != nil {
			return err
		}

		var query2 = bson.M{"$push": bson.M{"sections": s}}

		if err := User.UpdateQuery(selector, query2, false); err != nil {
			return err
		}
	}

	return nil
}
