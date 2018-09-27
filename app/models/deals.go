package models

import (
	"spyc_backend/app"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// Deal is the struct
type Deal struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	From        string        `json:"from" bson:"from"`
	To          bson.ObjectId `json:"to" bson:"to"`
	Amount      int64         `json:"amount" bson:"amount"`
	TotalAmount int64         `json:"current" bson:"current"`
	PublicK     string        `json:"public" bson:"public"`
	PrevTx      string        `json:"-" bson:"prevtx"`
}

// GetDocumentName ...
func (m *Deal) GetDocumentName() string {
	return "deals"
}

//FindTxs returns all transactions from User from DB ordered by creation date
func (m *Deal) FindTxs(userID bson.ObjectId) ([]Deal, error) {

	txs := []Deal{}
	Deal, _ := app.Mapper.GetModel(&Deal{})
	if err := Deal.Query(
		bson.M{"$and": []bson.M{
			bson.M{"from": userID}, bson.M{"to": userID},
		}}).Sort([]string{"created_at"}).Exec(&txs); err != nil {
		revel.ERROR.Print(err)
		return txs, err
	}
	return txs, nil
}

// AddAmount adds the given amount to the current transaction
func (m *Deal) AddAmount(amount int64) *Deal {
	m.Amount += amount
	return m
}

// AddPublicK adds the given key to the current transaction
func (m *Deal) AddPublicK(key string) *Deal {
	m.PublicK = key
	return m
}

// Save storage the current transaction into DB
func (m *Deal) Save() bool {
	Deal, _ := app.Mapper.GetModel(&Deal{})

	if err := Deal.Create(m); err != nil {
		revel.ERROR.Print(err)
		return false
	}
	return true
}

// AddCoinsFromSpychatter ...
func AddCoinsFromSpychatter(userID, amount, string int64) bool {
	tx := Deal{
		From:        "0",
		Amount:      amount,
		TotalAmount: amount,
		PrevTx:      "0",
	}
	if !tx.Save() {
		return false
	}
	return true
}
