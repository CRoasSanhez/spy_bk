package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/auth"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app/serializers"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// ChallengesController controller
type ChallengesController struct {
	BaseController
}

// CreateChallenge [/v2/challs] POST
// creates the challenge for the user
func (c ChallengesController) CreateChallenge() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var challenge models.Challenge
	if err := c.PermitParams(&challenge, false, "title", "secret_text", "properties", "type", "end_date", "players", "public"); err != nil {
		return c.ErrorResponse(nil, c.Message("error.badRequest"), core.ModelStatus[core.StatusInvalidID])
	}

	// If challenge not has players then return error
	if len(challenge.Players) == 0 {
		return c.ErrorResponse(nil, c.Message("error.empty", "Players"), core.ModelStatus[core.ValidationRequired])
	}

	// if Challenge, ok := app.Mapper.GetModel(&challenge); ok {
	challenge.UserID = c.CurrentUser.GetID()
	if c.CurrentUser.PersonalData.Address.Country == "" {
		challenge.Country = "N/A"
	} else {
		challenge.Country = c.CurrentUser.PersonalData.Address.Country
	}

	challenge.User = &struct {
		UserID bson.ObjectId `json:"user_id" bson:"user_id"`
		Email  string        `json:"email" bson:"email"`
	}{c.CurrentUser.GetID(), c.CurrentUser.Email}

	var completedUsers []models.PlayerInfo
	// Create Challengers with status and token init
	for _, uID := range challenge.Players {
		if bson.IsObjectIdHex(uID.UserID) && uID.UserID != c.CurrentUser.GetID().Hex() {
			completedUsers = append(completedUsers,
				models.PlayerInfo{
					UserID: uID.UserID,
					Status: core.StatusInit,
					Token:  core.StatusInit,
				})
		}
	}

	//challenge.Properties.Score = challenge.SecretText
	challenge.Players = completedUsers
	if challenge.Properties.Intents <= 0 {
		challenge.Properties.Intents = 1
	}

	if err := models.CreateDocument(&challenge); err != nil {
		revel.ERROR.Print("ERROR creating challenge --- " + err.Error())
		return c.ErrorResponse(nil, c.Message("error.create", "Challenge"), core.ValidationStatus[core.StatusError])
	}

	// Create Challenge Stats
	var fields = bson.M{"challenges.created": 1, "challenges.created_types." + challenge.Type: 1, "challenges.total": 1}
	if err := models.StatsUpdateFields(fields, []string{"challenges.last_access"}, c.CurrentUser.GetID()); err != nil {
		revel.ERROR.Print(err)
	}

	if c.Request.Locale == "" {
		c.Request.Locale = c.CurrentUser.Device.Language
	}

	if c.Request.Locale == "" {
		c.Request.Locale = "en"
	}

	// generate reward for the Challenge winner
	var reward = models.Reward{
		ResourceID:   challenge.GetID(),
		ResourceType: core.ModelTypeChallenge,
		Name:         map[string]interface{}{c.Params.Get("language"): "Challenge"},
		Description:  map[string]interface{}{c.Request.Locale: challenge.Title},
		Type:         core.ModelTypeChallenge,
		Coins:        200,
		Views:        0,
		Langs:        append([]string{}, c.Request.Locale),
		Multi:        false,
		MaxWinners:   1,
	}

	// if Reward, ok := app.Mapper.GetModel(&reward); ok {
	if err := models.CreateDocument(&reward); err != nil {
		revel.ERROR.Print("ERROR CREATE Reward for Challenge: " + challenge.GetID().Hex())
		return c.ServerErrorResponse()
	}

	var gameURL = core.GetGameBasePath() + ":webGame/?ich=:challengeID&actokn=:token&lng=:lang&tp=cht"
	gameURL = strings.Replace(gameURL, ":webGame", challenge.Type, 1)
	gameURL = strings.Replace(gameURL, ":challengeID", challenge.GetID().Hex(), 1)

	return c.SuccessResponse(models.ChallengeResponse{ID: challenge.GetID().Hex(), URL: gameURL}, "success", core.ModelsType[core.ModelSimpleResponse], nil)
}

// SaveChallengeFile [/v2/challs/:id/upload] POST
// attach a file(s) to challenge generating model reward
func (c ChallengesController) SaveChallengeFile(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		revel.ERROR.Print("ERROR invlid ID: " + id)
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var err error
	var challenge models.Challenge
	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {

		var match = bson.M{"$and": []bson.M{
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"user_id": c.CurrentUser.GetID()},
		}}
		var pipe = mgomap.Aggregate{}.LookUp("rewards", "_id", "resource_id", "rewards").Match(match)

		if err = Challenge.Pipe(pipe, &challenge); err != nil || len(challenge.Rewards) == 0 {
			revel.ERROR.Printf("Not found or NO REWARDS found for challenge: %s", challenge.GetID().Hex())
			return c.ErrorResponse(c.Message("error.notFound", "Challenge", ""), c.Message("error.notFound", "CHallenge", ""), 400)
		}

		if len(c.Params.Files["reward_file"]) != 1 {
			revel.ERROR.Print("ERROR no reward given")
			return c.ErrorResponse(nil, c.Message("error.empty", "File"), core.ModelStatus[core.StatusInvalidID])
		}

		var chllngFile models.ChallengeFiles

		// Upoload multipart file (normal reward image)
		if err = chllngFile.Attachment.Init(models.AsDocumentBase(&challenge), c.Params.Files["reward_file"][0]); err != nil {
			revel.ERROR.Print("ERROR INIT Reward Image --- " + err.Error())
			return c.ErrorResponse("", "Error Init image", 400)
		}

		if err = chllngFile.Attachment.Upload(); err != nil {
			revel.ERROR.Print("ERROR UPLOAD Reward Image --- " + err.Error())
			return c.ErrorResponse("", "Error uploading image", 400)
		}

		// Upload Image(byteArray) thumbnail
		var byteArray = c.ResizeImage(100, "reward_file", core.BlurChallenge)
		if len(byteArray) <= 0 {
			return c.ErrorResponse("", "Error Resizing image", 400)
		}

		if err = chllngFile.Thumbnail.UploadBytes(models.AsDocumentBase(&challenge), byteArray, c.Params.Files["reward_file"][0].Filename); err != nil {
			revel.ERROR.Print("ERROR UPLOAD Reward Thumbnail --- " + err.Error())
			return c.ErrorResponse("", "Error uploading Thumbnail", 400)
		}

		// Add last access stats
		if err = models.StatsNowFields(c.CurrentUser.GetID(), []string{"challenges.last_access"}); err != nil {
			revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", c.CurrentUser.GetID().Hex(), err.Error())
		}

		if Reward, ok := app.Mapper.GetModel(&models.Reward{}); ok {
			challenge.Rewards[0].IsVisible = challenge.IsPublic
			challenge.Rewards[0].Files = append(challenge.Rewards[0].Files, chllngFile)
			if err = Reward.Update(&challenge.Rewards[0]); err != nil {
				revel.ERROR.Printf("ERROR updating reward: %s --- %s", challenge.Rewards[0].GetID().Hex(), err.Error())
				return c.ServerErrorResponse()
			}
		}

		return c.SuccessResponse("success", "success", core.ModelsType[core.ModelSimpleResponse], nil)
	}

	return c.ServerErrorResponse()
}

// GetChallengeInfo [/v2/webgames/chat/:id] GET
// receives the user token and returns challenge info to WebGame
func (c ChallengesController) GetChallengeInfo(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		revel.ERROR.Print("ERROR invalid ID: " + id)
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var challenge models.Challenge
	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {
		var challenger models.PlayerInfo

		if err := Challenge.Find(id).Exec(&challenge); err != nil || challenge.Status.Name != core.StatusInit {
			revel.ERROR.Print("ERROR challenge not found: " + id)
			return c.ErrorResponse(nil, c.Message("error.notFound", "Challenge"), core.ValidationStatus[core.StatusNotFound])
		}

		// Generate new token with status active
		token, err := auth.GenerateChallengeToken(id, c.CurrentUser.GetID().Hex(), core.StatusActive)
		if err != nil {
			revel.ERROR.Print("ERROR generating token")
			return c.ErrorResponse(nil, c.Message("error.generateNew", "Token"), core.ValidationStatus[core.StatusError])
		}

		// Find Player with status and token init
		players := challenge.Players
		for i := 0; i < len(players); i++ {

			// Find current user in chaallenge players list
			if players[i].UserID == c.CurrentUser.GetID().Hex() && players[i].Status == core.StatusInit && players[i].Token == core.StatusInit {
				challenger.UserID = players[i].UserID
				challenger.UserName = c.CurrentUser.UserName
				challenger.Status = core.StatusActive
				players[i].Status = core.StatusActive
				players[i].Token = token
				players[i].UserName = c.CurrentUser.UserName
				break
			}
		}

		if challenger.UserID == "" {
			revel.ERROR.Print("ERROR REQUEST Player not found or doesn't allowed: " + c.CurrentUser.GetID())
			return c.ErrorResponse(nil, c.Message("error.notAllowed", "Challenge"), core.ValidationStatus[core.StatusError])
		}

		if err := Challenge.Update(&challenge); err != nil {
			revel.ERROR.Print("ERROR updating challenge")
			return c.ErrorResponse(nil, c.Message("error.update", "Challenge"), core.ValidationStatus[core.StatusError])
		}

		// Update User Stats
		if err = models.StatsNowFields(c.CurrentUser.GetID(), []string{"challenges.last_access"}); err != nil {
			revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", c.CurrentUser.GetID().Hex(), err.Error())
		}

		challenge.Player = challenger
		challenge.PlayerToken = token
		challenger.Token = token

		return c.SuccessResponse(challenge, "success", core.ModelsType[core.ModelTypeChallenge], serializers.ChallengeInfoSerializer{})
	}
	return c.ServerErrorResponse()
}

// StartChallenge [/v2/webgames/chat/start] POST
// changes challenger status to penidng validation to WebGame
func (c ChallengesController) StartChallenge() revel.Result {
	// Authenticate the game token with status active
	challenge, user, err := auth.AuthenticateChallengeToken(c.Request, core.StatusActive, core.StatusInit)
	if err != nil {
		revel.ERROR.Printf("Error Invalid Token -- %s", err.Error())
		return c.ForbiddenResponse()
	}

	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {
		var answerReq models.WebGameFinish
		var newToken string
		var isOK bool

		if err = c.PermitParams(&answerReq, true, "data", "status"); err != nil || answerReq.Status != core.ActionStart {
			return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
		}

		//Create the token with action start
		newToken, err = auth.GenerateChallengeToken(challenge.GetID().Hex(), user.GetID().Hex(), core.ActionStart)
		if err != nil {
			revel.ERROR.Printf("Error generating token --- %s", err.Error())
			return c.ErrorResponse(nil, c.Message("error.generateNew", "Token"), core.ValidationStatus[core.StatusError])
		}

		// Find Player with status and token init
		players := challenge.Players
		for i := 0; i < len(players); i++ {
			if players[i].UserID == user.GetID().Hex() && players[i].Status == core.StatusActive && players[i].Token != core.StatusInit {
				players[i].Status = core.ActionStart
				players[i].Token = newToken
				players[i].Score = answerReq.Data
				isOK = true
				break
			}
		}

		if !isOK {
			revel.ERROR.Print("ERROR REQUEST User not allowed or doesn't exists: " + user.GetID())
			return c.ErrorResponse(nil, c.Message("error.notAllowed", "Challenge"), core.ValidationStatus[core.StatusNotFound])
		}

		// Update User Stats
		if err = models.StatsNowFields(user.GetID(), []string{"challenges.last_access"}); err != nil {
			revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", user.GetID().Hex(), err.Error())
		}

		if err = Challenge.Update(&challenge); err == nil {
			return c.SuccessResponse(bson.M{"new_token": newToken}, "success", core.ModelsType[core.ModelWebGameStart], nil)
		}
		revel.ERROR.Print("ERROR updating challenge: " + challenge.GetID())
	}

	return c.ServerErrorResponse()
}

// ValidateChallenge [/v2/webgames/chat/validate] POST
// validates user's status in challenge and if won change it to complete
func (c ChallengesController) ValidateChallenge() revel.Result {

	challenge, user, err := auth.AuthenticateChallengeToken(c.Request, core.ActionStart, core.StatusInit)
	if err != nil {
		revel.ERROR.Printf("Error Invalid Token -- %s", err.Error())
		return c.ForbiddenResponse()
	}

	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {
		var answerReq models.WebGameFinish
		var isWinner bool
		var userExists bool

		// Bind response from request
		if err := c.PermitParams(&answerReq, true, "data", "status"); err != nil {
			return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
		}

		//Create the token with action upadte
		var newToken string
		newToken, err = auth.GenerateChallengeToken(challenge.GetID().Hex(), user.GetID().Hex(), core.ActionUpdate)
		if err != nil {
			revel.ERROR.Print("Error generating token")
			return c.ErrorResponse(nil, c.Message("error.generateNew", "Token"), core.ValidationStatus[core.StatusError])
		}

		var players = challenge.Players
		for i := 0; i < len(players); i++ {
			if players[i].UserID == user.GetID().Hex() && players[i].Status == core.ActionStart && players[i].Token != core.StatusInit {
				userExists = true
				players[i].Token = core.StatusCompleted
				players[i].CompletedAt = time.Now()
				players[i].Score = answerReq.Data

				if answerReq.Status != core.StatusSuccess {
					players[i].Status = core.StatusCompletedNotWinner
					break
				}

				// Validate challenge is hangman and the score given is the same as the stored
				if challenge.Type == "hangman" && answerReq.Data == players[i].Score {
					players[i].Status = core.StatusCompletedWinner
					isWinner = true
					break
				}

				// Validate is not type hangman
				if challenge.Type != "hangman" {
					players[i].Status = core.StatusCompletedWinner
					players[i].Token = newToken
					isWinner = true
				} else {
					players[i].Status = core.StatusCompletedNotWinner
				}
				break
			}
		}
		if !userExists {
			revel.ERROR.Print("ERROR REQUEST User not allowed or doesn't exists: " + user.GetID().Hex())
			return c.ErrorResponse(nil, c.Message("error.notAllowed", "Challenge"), core.ValidationStatus[core.StatusWrong])
		}

		// Updating last result date
		if err = models.StatsSetField(user.GetID(), bson.M{"challenges.last_result": answerReq.Status}); err != nil {
			revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", user.GetID().Hex(), err.Error())
		}

		// If player is NOT WINNER
		if !isWinner {
			if err = models.StatsUpdateFields(bson.M{"challenges.uncompleted": 1, "challenges.total": 1}, []string{"challenges.last_access", "challenges.last_uncompleted"}, c.CurrentUser.GetID()); err != nil {
				revel.ERROR.Printf("ERROR UPDATE STATS challenge: %s ", err.Error())
			}
			if err = Challenge.Update(&challenge); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}
			return c.SuccessResponse(bson.M{"new_token": ""}, "success", core.SubscriptionStatus[core.StatusCompletedNotWinner], nil)
		}

		challenge.SetStatus(core.StatusCompleted)

		if err = Challenge.Update(&challenge); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		// Update User Stats
		if err = models.StatsUpdateFields(bson.M{"challenges.completed": 1, "challenges.total": 1}, []string{"challenges.last_access", "challenges.last_completed"}, c.CurrentUser.GetID()); err != nil {
		}

		var reward models.Reward
		var params = []bson.M{
			bson.M{"status.name": core.StatusInit}, bson.M{"resource_id": challenge.GetID()}, bson.M{"type": core.ModelTypeChallenge},
		}

		if Reward, ok := app.Mapper.GetModel(&reward); ok {
			if err := Reward.FindWithOperator("$and", params).Exec(&reward); err != nil {
				revel.ERROR.Printf("ERROR FIND reward for challenge: %s --- %s", challenge.GetID().Hex(), err.Error())
				return c.ServerErrorResponse()
			}

			// Create winner
			win := models.Win{}
			_, errWinner := win.CreateWin(reward.GetID(), user.GetID(), core.ModelTypeChallenge, answerReq.Data)
			if errWinner != nil {
				return c.ErrorResponse(c.Message("error.create", "Winner"), c.Message("error.create", "Winner"), core.GameStatus[core.StatusInactive])
			}

			if reward.CompleteReward(user.GetID().Hex()) {
				return c.SuccessResponse(bson.M{"new_token": newToken}, "success", core.SubscriptionStatus[core.StatusCompletedWinner], nil)
			}
		}
	}
	return c.ServerErrorResponse()
}

// UpdateChallengeScore [/v2/webgames/chat] PUT
// updates the user score if winner
func (c ChallengesController) UpdateChallengeScore() revel.Result {
	challenge, user, err := auth.AuthenticateChallengeToken(c.Request, core.ActionUpdate, core.StatusCompleted)
	if err != nil {
		revel.ERROR.Printf("Error Invalid Token -- %s", err.Error())
		return c.ForbiddenResponse()
	}

	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {
		var answerReq models.WebGameFinish
		var userExists bool

		// Bind response from request
		if err := c.PermitParams(&answerReq, true, "data", "status"); err != nil {
			return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
		}

		players := challenge.Players
		for i := 0; i < len(players); i++ {
			if players[i].UserID == user.GetID().Hex() && players[i].Status == core.StatusCompletedWinner {
				userExists = true
				players[i].Score = answerReq.Data
				break
			}
		}
		if !userExists {
			revel.ERROR.Print("ERROR REQUEST User not allowed or doesn't exists: " + user.GetID().Hex())
			return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
		}

		// Update User Stats
		if err = models.StatsNowFields(user.GetID(), []string{"challenges.last_access"}); err != nil {
			revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", user.GetID().Hex(), err.Error())
		}

		if err = Challenge.Update(&challenge); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
		return c.SuccessResponse("success", "Challenge Updated", core.SubscriptionStatus[core.StatusCompletedWinner], nil)
	}

	return c.ServerErrorResponse()
}

// GetChallengeUserStatus [/v2/challs/:id/status] GET
// returns the user status in the challenge
func (c ChallengesController) GetChallengeUserStatus(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		revel.ERROR.Print("Error Invalid ID")
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var player, msg, code = GetPlayerStatus(id, c.CurrentUser.GetID().Hex(), true)

	switch code {
	case core.SubscriptionStatus[core.StatusCompletedWinner], core.SubscriptionStatus[core.StatusWinnerMaxViewsReached], core.SubscriptionStatus[core.StatusCompletedNotWinner]:
		return c.SuccessResponse(player, msg, code, serializers.PlayerSerializer{})
	case 1:
		return c.ErrorResponse(c.Message("error.badRequest"), c.Message("error.notFound", "", ""), 400)
	}
	return c.ServerErrorResponse()
}

// GetChallengePlayerStatus [/v2/challs/:id/status/:idPlayer] GET
// returns the status for the given player in the challenge
func (c ChallengesController) GetChallengePlayerStatus(id, idPlayer string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		revel.ERROR.Print("Error Invalid ID")
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var player models.PlayerInfo
	var msg string
	var code int

	if idPlayer == "" || idPlayer == c.CurrentUser.GetID().Hex() {
		player, msg, code = GetPlayerStatus(id, c.CurrentUser.GetID().Hex(), true)
	} else {
		if bson.IsObjectIdHex(idPlayer) {
			player, msg, code = GetPlayerStatus(id, idPlayer, false)
		} else {
			return c.ErrorResponse(c.Message("error.badRequest"), c.Message("error.badRequest"), 400)
		}
	}

	switch code {
	case core.SubscriptionStatus[core.StatusCompletedWinner], core.SubscriptionStatus[core.StatusWinnerMaxViewsReached], core.SubscriptionStatus[core.StatusCompletedNotWinner]:
		return c.SuccessResponse(player, msg, code, serializers.PlayerSerializer{})
	case 1:
		return c.ErrorResponse(c.Message("error.badRequest"), c.Message("error.notFound", "", ""), 400)
	}
	return c.ServerErrorResponse()
}

// GetChallengeStatus [/v2/challs/status/:id] GET
// returns the full challenge status available for the challenge owner
func (c ChallengesController) GetChallengeStatus(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		revel.ERROR.Print("Error Invalid ID")
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var params = []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"user_id": c.CurrentUser.GetID()},
	}
	var challenge models.Challenge
	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {
		if err := Challenge.FindWithOperator("$and", params).Exec(&challenge); err != nil {
			revel.ERROR.Printf("ERROR FIND Challenge: %s --- %s", id, err.Error())
			return c.ErrorResponse(nil, c.Message("error.notFound", "Challenge"), core.ModelStatus[core.StatusNotFound])
		}
		return c.SuccessResponse(challenge, "General Info", core.ModelsType[core.ModelTypeChallenge], serializers.ChallengeStatusSerializer{})
	}
	return c.ServerErrorResponse()
}

// NotifyChallengePlayers [/v2/challs/:id/notify] POST
// notifies the users in challenge through push notification
func (c ChallengesController) NotifyChallengePlayers(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var challenge models.Challenge
	match := bson.M{"$and": []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"user_id": c.CurrentUser.GetID()},
		bson.M{"status.name": core.StatusInit},
	}}

	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {

		pipe := mgomap.Aggregate{}.LookUp("rewards", "_id", "resource_id", "rewards").Match(match)

		if err := Challenge.Pipe(pipe, &challenge); err != nil || len(challenge.Rewards) == 0 {
			revel.ERROR.Printf("Not found or NO REWARDS found for challenge: %s", challenge.GetID().Hex())
			return c.ErrorResponse(c.Message("error.notFound", "CHallenge", ""), c.Message("error.notFound", "CHallenge", ""), 400)
		}

		go c.NotifyChallenge(challenge)
		return c.SuccessResponse("success", "success", core.ValidationStatus[core.StatusSuccess], nil)
	}

	return c.ServerErrorResponse()
}

// NotifyChallenge notifies users that they have a new challenge PUSH NOTIFICATION
func (c ChallengesController) NotifyChallenge(challenge models.Challenge) {

	var user = models.User{}
	var game = models.WebGame{}
	User, _ := app.Mapper.GetModel(&user)

	game = game.FindByNameURL(challenge.Type)
	if game.Name == "" {
		revel.ERROR.Print("Game with name_url: " + challenge.Type + " not found")
		return
	}

	// Find game to send itsimage in push notification
	var sender = &mgo.DBRef{
		Id:         c.CurrentUser.GetID(),
		Collection: c.CurrentUser.GetDocumentName(),
		Database:   app.Mapper.DatabaseName,
	}

	var senderName = c.CurrentUser.PersonalData.FirstName + " " + c.CurrentUser.PersonalData.LastName
	var title = "Complete the game " + challenge.Title

	for _, player := range challenge.Players {
		if err := User.Find(player.UserID).Exec(&user); err == nil {
			var receiver = &mgo.DBRef{
				Id:         user.GetID(),
				Collection: user.GetDocumentName(),
				Database:   app.Mapper.DatabaseName,
			}

			invite := models.Invitation{
				Sender:       sender,
				Receiver:     receiver,
				Resource:     sender,
				ResourceType: challenge.GetDocumentName(),
				Type:         "challenge",
			}

			if Invite, ok := app.Mapper.GetModel(&models.Invitation{}); ok {
				if err := Invite.Create(&invite); err != nil {
					revel.ERROR.Print(err)
				}
			}
			go c.NewNotification("show", game.Attachment, senderName+"has challenged you", "challenge", title, "invite", user, challenge.GetID(), true)
		}
	}
}

// Delete [/challs] DELETE
// Deletes the challenge if the user is the owner
func (c ChallengesController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		revel.ERROR.Print("ERROR invlid ID: " + id)
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	if Challenge, ok := app.Mapper.GetModel(&models.Challenge{}); ok {

		var selector = bson.M{"user_id": c.CurrentUser.GetID(), "_id": bson.ObjectIdHex(id)}
		var query = bson.M{"$set": bson.M{"deleted": true, "status.name": core.StatusInactive}}

		if err := Challenge.UpdateQuery(selector, query, false); err != nil {
			revel.ERROR.Printf("ERROR DELETE Challenge %s --- %s", id, err.Error())
			return c.ErrorResponse(err, err.Error(), 400)
		}
		return c.SuccessResponse("success", "Deleted", 200, nil)
	}
	return c.ServerErrorResponse()
}

// GetPlayerStatus returns the player info according to the challenge type (public/private)
func GetPlayerStatus(challengeID, userID string, current bool) (player models.PlayerInfo, msg string, code int) {
	var challenge models.Challenge
	code, msg = core.SubscriptionStatus[core.StatusCompletedNotWinner], "Player Info"

	match := bson.M{"$and": []bson.M{bson.M{"_id": bson.ObjectIdHex(challengeID)}}}
	pipe := mgomap.Aggregate{}.LookUp("rewards", "_id", "resource_id", "rewards").Match(match)

	if Challenge, ok := app.Mapper.GetModel(&challenge); ok {

		if err := Challenge.Pipe(pipe, &challenge); err != nil || len(challenge.Rewards) == 0 {
			revel.ERROR.Printf("Not found or NO REWARDS found for challenge: %s", challenge.GetID().Hex())
			return player, "", 1
		}

		player = FindPlayer(userID, challenge.Players)

		// Return code error if challenge is not public or no player found
		if !challenge.IsPublic || player.UserID == "" {
			return player, "", 1
		}

		// Set user stats
		if err := models.StatsNowFields(bson.ObjectIdHex(userID), []string{"challenges.last_access"}); err != nil {
			revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", userID, err.Error())
		}

		var reward = challenge.Rewards[0]
		player.Goal = challenge.Properties.Score
		player.Game = challenge.Type

		if player.Status == core.StatusCompletedWinner {
			code, msg = core.SubscriptionStatus[core.StatusCompletedWinner], "Winner"
		}

		if Reward, ok := app.Mapper.GetModel(&reward); ok {

			if current && player.Status == core.StatusCompletedWinner {
				var selector = bson.M{"$and": []bson.M{
					bson.M{"$or": []bson.M{
						bson.M{"user_id": userID},
						bson.M{"users": bson.M{"$elemMatch": bson.M{"$eq": userID}}},
					}},
					bson.M{"resource_id": challenge.GetID()},
					bson.M{"resource_type": core.ModelTypeChallenge},
				}}

				var query = bson.M{"$set": bson.M{"views": (reward.Views + 1), "status.name": core.StatusCompleted}}
				if reward.Views >= core.MaxRewardViews {
					query = bson.M{"$set": bson.M{"views": (reward.Views + 1), "status.name": core.StatusObtained}}
					reward.Files = RemoveFromReward("file", reward.Files)
					code, msg = core.SubscriptionStatus[core.StatusWinnerMaxViewsReached], "Winner but Max views Reached"
				}

				if err := Reward.UpdateQuery(selector, query, false); err != nil {
					revel.ERROR.Printf("ERROR UPDATE Reward views: %s --- %s", reward.GetID().Hex(), err.Error())
					return player, "", 0
				}
			} else {
				reward.Files = RemoveFromReward("file", reward.Files)
			}
			player.Rewards = append(player.Rewards, reward)
			return player, msg, code
		}
	}
	return player, "", 0
}

// FindPlayer ...
func FindPlayer(key string, array []models.PlayerInfo) models.PlayerInfo {
	for i := 0; i < len(array); i++ {
		if array[i].UserID == key {
			return array[i]
		}
	}
	return models.PlayerInfo{}
}

// RemoveFromReward ...
func RemoveFromReward(fileType string, array []models.ChallengeFiles) []models.ChallengeFiles {
	for i := 0; i < len(array); i++ {
		if fileType == "file" {
			array[i].Attachment = models.Attachment{}
		} else {
			array[i].Thumbnail = models.Attachment{}
		}
	}
	return array
}
