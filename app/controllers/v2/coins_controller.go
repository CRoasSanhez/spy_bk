package v2

import (
	"github.com/revel/revel"

	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
)

// CoinsController ...
type CoinsController struct {
	BaseController
}

// SendTransaction ...
func (c CoinsController) SendTransaction() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// Values for channel function
	var prevID string
	var amountFrom int64
	//var amountTo int64

	// Bind json request
	var req core.SimpleTransaction

	if err := c.Params.BindJSON(&req); err != nil {
		revel.ERROR.Print(err)
		return c.ErrorResponse(nil, "Bad Request", core.ValidationStatus[core.StatusError])
	}

	// Validate amount is not greater than the allowed
	if req.Amount > core.MaxTransactionAmount {
		return c.ErrorResponse(nil, "Invalid coins", core.ValidationStatus[core.StatusError])
	}

	chanPrevID := make(chan string)
	chanAmountFrom := make(chan int64)
	//chanAmountTo := make(chan int64)

	go GetUserTransactions(c.CurrentUser.GetID(), chanAmountFrom, chanPrevID)
	//go core.GetUserTransactions(bson.ObjectIdHex(req.To), chanAmountTo,chanPrevID)

	prevID = <-chanPrevID
	amountFrom = <-chanAmountFrom
	//amountTo = <-chanAmountTo

	// Verify there is no validation error and the sender has more coins than the amount to send
	if amountFrom == 0 || req.Amount > amountFrom {
		return c.ErrorResponse(nil, "Invalid Transaction", core.ValidationStatus[core.StatusError])
	}

	tx := models.Deal{From: c.CurrentUser.GetID().Hex(), To: bson.ObjectIdHex(req.To), PrevTx: prevID}

	if !tx.AddAmount(req.Amount).Save() {
		return c.ErrorResponse(nil, "Error saving transaction", core.ValidationStatus[core.StatusError])
	}

	// Update sender transaction in User Model
	c.CurrentUser.AddTransaction(req.Amount, "add")

	if User, ok := app.Mapper.GetModel(&c.CurrentUser); ok {
		if err := User.Find(req.To).Exec(&c.CurrentUser); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		// Update receiver transaction in user Model
		c.CurrentUser.AddTransaction(req.Amount, "remove")

		return c.SuccessResponse("success", "success", core.ModelsType[core.ModelSimpleResponse], nil)

		//---------------------------------
		// BITCOIN implementation
		/*
			if from == "" || to == "" || payload == "" {
				return core.ResponseERR("Incomplete fields")
			}
		*/

		/*
			rxCoins := &models.RxCoins{}

			err := binding.Bind(c.Request.Request, rxCoins)
			if err != nil {
				return core.ResponseERR(err.Error())
			}

			byteFrom := []byte("cesar")
			byteTo := []byte("god")
			bytePayload := []byte("500")

			log.Printf("From: %s\n", byteFrom)
			log.Printf("To: %s\n", byteTo)
			log.Printf("Payload: %s\n", bytePayload)

			t := blockchain.NewTransaction(byteFrom, byteTo, bytePayload)

			return core.ResponseJSON{Data: t, Status: http.StatusOK}
		*/
	}

	return c.ServerErrorResponse()
}

// GetUserTransactions ...
func GetUserTransactions(userID bson.ObjectId, totalAmount chan int64, prevID chan string) {
	// validate sender has the specific amount
	userTxs := []models.Deal{}
	tx := models.Deal{}
	var ok = false
	var totalSum int64
	totalSum = 0

	//Deal,_:= app.Mapper.GetModel(&tx)
	userTxs, _ = tx.FindTxs(userID)

	for _, t := range userTxs {
		ok = false
		if t.To == userID {
			totalSum += t.Amount
			ok = true
			continue
		}
		if bson.ObjectIdHex(t.From) == userID {

			// Validate that transaction was successfully and matches the totalSum
			if t.TotalAmount >= t.Amount && totalSum == t.TotalAmount {
				totalSum -= t.Amount
				ok = true
				continue
			}
		}
		if ok == false {
			totalSum = 0
			break
		}
	}

	// Validate user has enough coins
	prevID <- userTxs[len(userTxs)-1].GetID().Hex()
	totalAmount <- totalSum
}
