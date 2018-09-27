package models

import (
	"errors"
	"math"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"time"

	"gopkg.in/mgo.v2/bson"

	"golang.org/x/crypto/bcrypt"

	"github.com/Reti500/mgomap"
	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/is"
)

// UserHistorial ...
type UserHistorial struct {
	Amount      int64 `json:"amount" bson:"amount"`
	TotalAmount int64 `json:"current" bson:"current"`
}

// User ...
type User struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Email        string          `json:"email" bson:"email"`
	UserName     string          `json:"user_name" bson:"user_name"`
	PasswordHash string          `json:"-" bson:"password_hash"`
	Friends      []string        `json:"-" bson:"friends"`
	BlackList    []string        `json:"-"  bson:"blacklist"`
	PersonalData PersonalData    `json:"personal_data" bson:"personal_data"`
	Geolocation  *Geo            `json:"geolocation" bson:"geolocation"`
	Device       Device          `json:"device" bson:"device"`
	Saas         Saas            `json:"-" bson:"saas"`
	Attachment   Attachment      `json:"-" bson:"attachment"`
	Permissions  []Permission    `json:"-" bson:"permissions"`
	Historial    []UserHistorial `json:"-" bson:"historial"`
	Sections     []Section       `json:"-" bson:"sections"`

	// Extra parameters
	ExtraParameters map[string]interface{} `bson:"extra_parameters"`

	// Not saved fields
	Password string `json:"password" bson:"-"` // Not saved field
}

// Validate ...
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Geolocation, validation.Required),
	)
}

// GetDocumentName ...
func (u *User) GetDocumentName() string {
	return "users"
}

// GetPublicFields retuerns only specific fields form user
func (u User) GetPublicFields() []string {
	return []string{"user_name", "personal_data", "saas.name", "attachment"}
}

// PublicProfile returns only public fields
func (u *User) PublicProfile(id string, fields ...string) error {
	var f []string

	if !bson.IsObjectIdHex(id) {
		return ErrID
	}

	if len(fields) <= 0 {
		f = u.GetPublicFields()
	} else {
		f = fields
	}

	if Document, ok := app.Mapper.GetModel(u); ok {
		return Document.Find(id).Select(f).Exec(u)
	}

	return ErrDocument
}

// FullName ...
func (u User) FullName() string {
	return u.PersonalData.FirstName + " " + u.PersonalData.LastName
}

// GeneratePassword ...
func (u *User) GeneratePassword() error {
	if u.Password == "" {
		return errors.New("Invalid password")
	}

	hash, err := MD5Crypt(u.Password)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)

	return nil
}

// SetState ...
func (u *User) SetState(state string) {
	u.Status.Name = state
	u.Status.Code = core.AccountStatus[state]
}

// MD5Crypt ...
func MD5Crypt(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// RemoveCoins removes coins from the user struct
func (uh *UserHistorial) RemoveCoins(coins int64) {
	uh.TotalAmount -= coins
	uh.Amount = -coins
}

// AddCoins adds coins to the user struct
func (uh *UserHistorial) AddCoins(coins int64) {
	uh.TotalAmount += coins
	uh.Amount = coins
}

// AddTransaction ...
func (u *User) AddTransaction(coins int64, transType string) {
	uh := UserHistorial{TotalAmount: u.Historial[len(u.Historial)-1].TotalAmount}

	if transType == "add" {
		uh.AddCoins(coins)
	} else {
		uh.RemoveCoins(coins)
	}
	u.Historial = append(u.Historial, uh)

	User, _ := app.Mapper.GetModel(&User{})
	User.Update(u)
}

// CanManageSection return true if current user can access to seccion
func (u *User) CanManageSection(section string) bool {
	for _, v := range u.Permissions {
		if v.Section == section && v.Action == core.ActionManage {
			return true
		}
	}

	return false
}

// AddFriend addiding friend id to friends list
func (u *User) AddFriend(friendID string) error {
	if User, ok := app.Mapper.GetModel(u); ok {
		var selector = bson.M{"_id": u.GetID()}
		var query = bson.M{"$addToSet": bson.M{"friends": friendID}}

		if err := User.UpdateQuery(selector, query, false); err != nil {
			return err
		}
	}

	return nil
}

// GetAge returns age of current user
func (u User) GetAge() int {
	var t = time.Now().Sub(u.PersonalData.BirthDate.Time).Hours()
	years, _ := math.Modf(t / 8760)

	return int(years)
}
