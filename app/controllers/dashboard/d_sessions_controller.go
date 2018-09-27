package dashboard

import (
	"net/http"
	"spyc_backend/app"
	"spyc_backend/app/auth"
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

type DSessionsController struct {
	DBaseController
}

func (c DSessionsController) New() revel.Result {
	return c.Render()
}

func (c DSessionsController) Create() revel.Result {
	var user models.User

	if User, ok := app.Mapper.GetModel(&user); ok {
		var pipe = mgomap.Aggregate{}.Match(bson.M{
			"email": c.Params.Get("email"),
			"permissions": bson.M{
				"$elemMatch": bson.M{
					"action":  core.ActionManage,
					"section": core.SectionDashboard,
				},
			},
		})

		if err := User.Pipe(pipe, &user); err != nil {
			return c.Redirect("/spyc_admin/sessions/new")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(c.Params.Get("password"))); err != nil {
			return c.Redirect("/spyc_admin/sessions/new")
		}

		token, err := auth.GenerateToken(user.GetID().Hex(), core.ActionAuth)
		if err != nil {
			return c.Redirect("/spyc_admin/sessions/new")
		}

		c.Session["access_token"] = token
		atoken := &http.Cookie{Name: "atoken", Value: token}
		c.SetCookie(atoken)

		return c.Redirect("/spyc_admin")
	}

	return c.ServerErrorResponse()
}

func (c DSessionsController) Delete() revel.Result {
	c.Session["access_token"] = ""

	return c.Redirect("/spyc_admin/sessions/new")
}
