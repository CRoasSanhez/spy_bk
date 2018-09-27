package dashboard

import (
	"math"
	"spyc_backend/app"
	"spyc_backend/app/auth"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// DUsersController controller
type DUsersController struct {
	DBaseController
}

// Index ....
func (c DUsersController) Index(page int) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	if page == 0 {
		page = 1
	}

	var users []models.User

	if User, ok := app.Mapper.GetModel(&models.User{}); ok {
		if err := User.All().Sort([]string{"email"}).Paginate((page - 1), 10).Limit(10).Exec(&users); err != nil {
			c.Flash.Error(err.Error(), err)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(users, "success", 200, serializers.UserSerializer{})
		}

		if count, err := User.All().Count(); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		} else {
			c.ViewArgs["Users"] = users
			c.ViewArgs["CurrentPage"] = page - 1
			c.ViewArgs["Pages"] = int(math.Ceil(float64(count) / float64(10)))
		}

		return c.Render()
	}

	return c.ServerErrorResponse()
}

// Create ...
func (c DUsersController) Create(role string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var user models.User

	if err := c.Params.BindJSON(&user); err != nil {
		revel.ERROR.Print(err)
		return c.ErrorResponse(err, err.Error(), 400)
	}

	user.Status = mgomap.Status{}
	user.SetState(core.StatusCreated)

	// Validations
	if err := user.Validate(); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	// Generamos la instancia para la conexion a la base de datos.
	if User, ok := app.Mapper.GetModel(&user); ok {
		// Generamos el password
		user.GeneratePassword()

		if err := User.Create(&user); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		t, err := auth.GenerateToken(user.GetID().Hex(), core.ActionAuth)
		if err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(user, "Created!!", 200, serializers.UserSerializer{Token: t})
	}

	return c.ServerErrorResponse()
}

// Show ...
func (c DUsersController) Show(id string) revel.Result {
	user := models.User{}
	User, _ := app.Mapper.GetModel(&user)

	err := User.Find(id).Exec(&user)
	if err != nil {
		c.Flash.Error(err.Error())
	}

	c.ViewArgs["Current"] = user

	return c.Render()
}

// ResetPassword ...
func (c DUsersController) ResetPassword() revel.Result {
	token := c.Params.Get("token")

	var t models.Token
	if err := models.GetDocumentBy(bson.M{"text": token}, &t); err != nil {
		return c.NotFound("Page not found")
	}

	c.Session["reset_pass"] = token

	return c.Render()
}

// ResetPasswordSuccess ...
func (c DUsersController) ResetPasswordSuccess() revel.Result {
	return c.Render()
}

// ChangePassword ...
func (c DUsersController) ChangePassword(password, confirmation string) revel.Result {
	if password != confirmation {
		return c.Redirect(DUsersController.ResetPassword)
	}

	var t models.Token
	var query = bson.M{"text": c.Session["reset_pass"]}
	if err := models.GetDocumentBy(query, &t); err != nil {
		return c.NotFound("Page not found")
	}

	// claims, err := auth.VerifyToken(c.Session["reset_pass"], core.ActionResetPass)
	// if err != nil {
	// 	return c.NotFound("Page not found")
	// }
	//
	// var user models.User
	//
	// if User, ok := app.Mapper.GetModel(&user); ok {
	// 	if err := User.Find(claims.ID).Exec(&user); err != nil {
	// 		return c.RenderError(err)
	// 	}
	//
	// 	user.Password = password
	// 	user.GeneratePassword()
	//
	// 	if err := User.Update(&user); err != nil {
	// 		return c.Redirect(DUsersController.ResetPassword)
	// 	}
	// }

	return c.Redirect("/spyc_admin/password_success")
}
