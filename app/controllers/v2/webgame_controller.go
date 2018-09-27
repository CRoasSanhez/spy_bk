package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/auth"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// WebGamesController handle requests from HTML Games
type WebGamesController struct {
	BaseController
}

// StartWebGame [/v2/webgames/cshnt/start] POST
// initializes the timer of webGame and creates new token
func (c WebGamesController) StartWebGame() revel.Result {
	// Authenticate the game token with status active
	targetSubsciption, user, err := auth.AuthenticateGameToken(c.Request, core.StatusActive)
	if err != nil {
		revel.ERROR.Printf("Error Invalid Token -- %s", err.Error())
		return c.ForbiddenResponse()
	}

	var newToken string

	if TargetUser, ok := app.Mapper.GetModel(&targetSubsciption); ok {
		var answerReq models.WebGameFinish

		// Bind response from request
		if err = c.PermitParams(&answerReq, false, "data", "status"); err != nil || answerReq.Status != core.ActionStart {
			revel.ERROR.Printf("Error request JSON --- %s", err.Error())
			return c.ErrorResponse(nil, c.Message("error.badRequest"), core.ValidationStatus[core.StatusWrong])
		}

		// Create the token with action start
		newToken, err = auth.GenerateGameToken(targetSubsciption.GetID().Hex(), user.GetID().Hex(), core.ActionStart)
		if err != nil {
			return c.ErrorResponse(nil, c.Message("error.generateNew", "Token"), core.ValidationStatus[core.StatusError])
		}

		targetSubsciption.Token = newToken
		targetSubsciption.Intents = 1
		targetSubsciption.Score = answerReq.Data
		startDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
		targetSubsciption.StartDate = startDate

		if err = TargetUser.Update(&targetSubsciption); err != nil {
			return c.ErrorResponse(nil, c.Message("error.update", "Token"), core.ValidationStatus[core.StatusError])
		}

		if err = models.StatsUpdateFields(nil, []string{"missions.last_access"}, user.GetID()); err != nil {
		}

		return c.SuccessResponse(bson.M{"new_token": newToken}, "success", core.ModelsType[core.ModelWebGameStart], nil)
	}
	return c.ServerErrorResponse()
}

// ValidateWebGame [/v2/webgames/cshnt/validate] POST
// saves the user score/answer and returns new token if winner
func (c WebGamesController) ValidateWebGame() revel.Result {
	targetSubsciption, user, err := auth.AuthenticateGameToken(c.Request, core.ActionStart)
	if err != nil {
		revel.ERROR.Printf("Error Invalid Token -- %s", err.Error())
		return c.ForbiddenResponse()
	}

	var isWinner bool
	var mission models.Mission
	var missionSubscription models.MissionUser
	if Mission, ok := app.Mapper.GetModel(&mission); ok {

		TargetUser, ok2 := app.Mapper.GetModel(&targetSubsciption)
		if !ok2 {
			return c.ServerErrorResponse()
		}

		var match = bson.M{"$and": []bson.M{
			bson.M{"status.name": core.StatusActive},
			bson.M{"targets._id": targetSubsciption.TargetID},
			bson.M{"targets.status.name": core.StatusActive},
		}}
		var pipe = mgomap.Aggregate{}.LookUp("targets", "_id", "missions", "targets").Match(match)

		// Find Mission and its targets
		if err = Mission.Pipe(pipe, &mission); err != nil {
			revel.ERROR.Printf("ERROR PIPE Mission with Target: %s --- %s", targetSubsciption.TargetID.Hex(), err.Error())
			return c.ErrorResponse(nil, c.Message("error.badRequest"), core.ValidationStatus[core.StatusNotFound])
		}

		// Find target in nested targets in Mission by bsonID from target subscription
		var target models.Target
		if target = mission.FindTargetByID(targetSubsciption.TargetID); target.Type != "" {
			var answerReq models.WebGameFinish
			if err := c.PermitParams(&answerReq, false, "data", "status"); err != nil {
				return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
			}

			targetSubsciption.Token = ""

			// Validate target is webgame
			if target.Type == core.GameTypeGame && answerReq.Status == core.StatusSuccess {
				var gameType = strings.Split(target.WebURL, "/")[4]

				if gameType == "hangman" && targetSubsciption.Score == answerReq.Data {
					isWinner = true
				}

				if gameType != "hangman" {
					isWinner = true
				}
			}

			// Updating last result date
			if err = models.StatsSetField(user.GetID(), bson.M{"missions.last_result": answerReq.Status}); err != nil {
				revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", user.GetID().Hex(), err.Error())
			}

			// Change subscription status if NOT WINNER
			if !isWinner {

				if err = models.StatsUpdateFields(bson.M{"missions.uncompleted": 1}, []string{"missions.last_uncompleted", "missions.last_access"}, c.CurrentUser.GetID()); err != nil {
				}

				if err = targetSubsciption.CompleteSubscription(core.StatusCompletedNotWinner); err != nil {
					revel.ERROR.Printf("ERROR UPDATE TargetSubscription: %s --- %s", targetSubsciption.GetID().Hex(), err.Error())
					return c.ErrorResponse(nil, c.Message("error.update", "Target Subscription"), core.ValidationStatus[core.StatusError])
				}
				// Change missionSubscription to complete
				if err = missionSubscription.CompleteSubscriptionByMissionUser(target.Missions.Hex(), user.GetID().Hex(), core.StatusCompletedNotWinner); err != nil {
					return c.ErrorResponse(nil, c.Message("error.update", "Mission Subscription"), core.ValidationStatus[core.StatusError])
				}
				return c.SuccessResponse(bson.M{"new_token": ""}, "success", core.SubscriptionStatus[core.StatusCompletedNotWinner], nil)

			}

			// If WINNER update subscriptions and return success
			if err = models.StatsUpdateFields(bson.M{"missions.completed": 1, "missions.total": 1}, []string{"missions.last_access", "missions.last_completed"}, c.CurrentUser.GetID()); err != nil {
			}

			if err = targetSubsciption.CompleteSubscription(core.StatusCompleted); err != nil {
				revel.ERROR.Printf("ERROR UPDATE TargetSubscription: %s --- %s", targetSubsciption.GetID().Hex(), err.Error())
				return c.ErrorResponse(nil, c.Message("error.update", "Target Subscription"), core.ValidationStatus[core.StatusError])
			}

			// Create the token with action Update
			targetSubsciption.Token, err = auth.GenerateGameToken(targetSubsciption.GetID().Hex(), user.GetID().Hex(), core.ActionUpdate)
			if err != nil {
				return c.ErrorResponse(nil, c.Message("error.generateNew", "Token"), core.ValidationStatus[core.StatusError])
			}

			// If there's another target then return token
			if target.NextTargetID != "" {
				if isWinner {
					return c.SuccessResponse(bson.M{"new_token": targetSubsciption.Token}, "success", core.ValidationStatus[core.StatusCompletedWinner], nil)
				}
				return c.SuccessResponse(bson.M{"new_token": ""}, "success", core.SubscriptionStatus[core.StatusCompletedNotWinner], nil)
			}

			var reward models.Reward
			rewards := reward.GetRewardsByResource(target.GetID().Hex())
			if len(rewards) <= 0 {
				return c.SuccessResponse("No reward collected", "failed", core.SubscriptionStatus[core.StatusCompletedNotWinner], nil)
			}

			reward = rewards[0]

			if err = TargetUser.UpdateQuery(bson.M{"_id": targetSubsciption.GetID()}, bson.M{"$set": bson.M{"reward_id": reward.GetID()}}, false); err != nil {
				revel.ERROR.Printf("ERROR UPDATE Target Subscription: %s, Reward: %s --- %s", targetSubsciption.GetID(), reward.GetID(), err.Error())
			}

			// Change missionSubscription to complete
			if err = missionSubscription.CompleteSubscriptionByMissionUser(target.Missions.Hex(), user.GetID().Hex(), core.StatusCompletedWinner); err != nil {
				return c.ErrorResponse(nil, c.Message("error.update", "Mission Subscription"), core.ValidationStatus[core.StatusError])
			}

			if ok := reward.CompleteReward(user.GetID().Hex()); !ok {
				return c.ErrorResponse(nil, c.Message("error.update", "Reward"), core.GameStatus[core.StatusInactive])
			}

			if (len(rewards) == 1 && !reward.Multi) || (reward.Multi && reward.Status.Name == core.StatusCompleted) {
				if mission.Complete() {
					c.NotifyMissionComplete(mission)
				}
			}

			// Notify to Slack that the user has obtained the reward
			core.Notify(core.NewEntry, "Reward Obtained", "A Reward was Obtained", mission.Title.GetString(user.Device.Language),
				core.GetDashboardPath()+"users/"+user.GetID().Hex(), "", mission.Description.GetString(user.Device.Language), []interface{}{
					bson.M{"title": mission.Type, "value": core.ConcatArray(mission.Countries), "short": true},
				})

			// Create winner for the target
			var win models.Win
			if _, err = win.CreateWin(reward.GetID(), user.GetID(), core.WinTypeCashHunt, answerReq.Data); err != nil {
				return c.ErrorResponse(nil, c.Message("error.create", "Win"), core.GameStatus[core.StatusInactive])
			}

			return c.SuccessResponse(bson.M{"new_token": targetSubsciption.Token}, "success", core.SubscriptionStatus[core.StatusCompletedWinner], nil)

		}
	}

	return c.ServerErrorResponse()
}

// UpdateWebGameScore [/v2/webgames/cshnt] PUT
// updates the user score if winner
func (c WebGamesController) UpdateWebGameScore() revel.Result {
	targetSubsciption, user, err := auth.AuthenticateGameToken(c.Request, core.ActionUpdate)
	if err != nil {
		revel.ERROR.Printf("Error Invalid Token -- %s", err.Error())
		return c.ForbiddenResponse()
	}

	if targetSubsciption.UserID != user.GetID() {
		return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
	}

	var answerReq models.WebGameFinish

	if err = c.PermitParams(&answerReq, false, "data", "status"); err != nil {
		return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
	}

	targetSubsciption.Score = answerReq.Data

	if err = models.StatsUpdateFields(nil, []string{"missions.last_access"}, c.CurrentUser.GetID()); err != nil {
	}

	if err = models.UpdateDocument(&targetSubsciption); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("success", "Game Updated", core.SubscriptionStatus[core.StatusCompletedWinner], nil)
}

// GetWebGameStatus [/v2/webgames/:idTarget/status] GET
// returns a simple response whether a target has been won or not
func (c WebGamesController) GetWebGameStatus(idTarget string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// Validate ID valid
	if !bson.IsObjectIdHex(idTarget) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var target models.Target

	var query = bson.M{"$and": []bson.M{
		{"status.name": core.StatusActive},
		{"_id": bson.ObjectIdHex(idTarget)},
	}}

	if err := models.GetDocumentBy(query, &target); err != nil {
		revel.ERROR.Printf("ERROR FIND Target: %s --- %s", idTarget, err.Error())
		return c.ErrorResponse(err, err.Error(), 400)
	}

	var reward models.Reward
	var rewards []models.Reward

	rewards = reward.GetRewardsByUserAndResourceType(c.CurrentUser.GetID().Hex(), []string{core.GameTypeGame}, true, 0)

	for _, v := range rewards {
		if v.ResourceID.Hex() == idTarget {
			reward = v
		}
	}

	if err := models.StatsUpdateFields(nil, []string{"missions.last_access"}, c.CurrentUser.GetID()); err != nil {
	}

	if reward.GetID().Hex() == "" {
		return c.SuccessResponse(bson.M{"data": "No reward collected"}, "failed", core.SubscriptionStatus[core.StatusCompletedNotWinner], nil)
	}

	return c.SuccessResponse(reward, "Success Reward", core.ModelsType[core.ModelReward], serializers.RewardSerializer{})
}

// GetWebGames [/v2/webgames] GET
// returns a list of all active webgames
func (c WebGamesController) GetWebGames() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var webGames models.WebGames
	var query = bson.M{"$and": []bson.M{
		{"status.name": core.StatusActive},
	}}

	if err := models.GetDocumentBy(query, &webGames); err != nil {
		return c.ErrorResponse(nil, c.Message("error.notFound", "Games", ""), core.ValidationStatus[core.StatusError])
	}

	return c.SuccessResponse(webGames, "Success", core.ModelsType[core.ModelWebGame], serializers.WebGamesSerializer{Lang: c.Request.Locale})
}
