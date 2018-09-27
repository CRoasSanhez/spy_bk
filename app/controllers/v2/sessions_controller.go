package v2

import (
	"spyc_backend/app/auth"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"golang.org/x/crypto/bcrypt"

	"errors"

	"github.com/revel/revel"
)

// SessionsController controller
type SessionsController struct {
	BaseController
}

// Create ...
func (c SessionsController) Create(account, password string) revel.Result {
	if account == "" || password == "" {
		return c.ErrorResponse(nil, "Missing params", 400)
	}

	account = strings.ToLower(account)

	var accountErr = errors.New("Account and password not match")
	var user = models.User{}

	var query = bson.M{
		"$or": []bson.M{
			bson.M{"email": account},
			bson.M{"user_name": account},
		},
	}

	if err := models.GetDocumentBy(query, &user); err != nil {
		return c.ErrorResponse(accountErr, accountErr.Error(), 400)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return c.ErrorResponse(accountErr, accountErr.Error(), 400)
	}

	if token, err := auth.GenerateToken(user.GetID().Hex(), "auth"); err == nil {
		return c.SuccessResponse(user, "Created!!", 200, serializers.UserSerializer{Token: token})
	}

	return c.ErrorResponse(accountErr, accountErr.Error(), 400)
}
