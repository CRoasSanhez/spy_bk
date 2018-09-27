package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

type RewardsController struct {
	BaseController
}

// CollectReward [/v2/rewards/collect/:idResource] POST
// changes the win status status to complete and creates a win
func (c RewardsController) CollectReward(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var selector = []bson.M{
		bson.M{"user_id": c.CurrentUser.GetID().Hex()},
		bson.M{"_id": id},
		bson.M{"multi": false},
	}
	var query = bson.M{"$set": []bson.M{
		bson.M{"status.name": core.StatusObtained},
		bson.M{"status.code": core.ValidationStatus[core.StatusObtained]},
	}}

	// Get pending Rewards for the user
	if Reward, ok := app.Mapper.GetModel(&models.Reward{}); ok {
		if err := Reward.UpdateQuery(selector, query, false); err != nil {
			revel.ERROR.Print("ERROR Find")
			return c.ErrorResponse(err, err.Error(), 400)
		}
		return c.SuccessResponse(bson.M{"data": "Reward collected successfully"}, "success", core.ModelsType[core.ModelSimpleResponse], nil)
	}

	return c.ServerErrorResponse()
}

// GetRewards [/v2/rewards/:page] GET
// returns the 20 rewards for the user according to the given page
func (c RewardsController) GetRewards(page int) revel.Result {

	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	//ChangeRewardsModel() // Remove when finish production

	var reward models.Reward
	if Reward, ok := app.Mapper.GetModel(&reward); ok {
		var rewards = []models.Reward{}
		var match = bson.M{"$and": []bson.M{
			bson.M{"$or": []bson.M{
				bson.M{"user_id": c.CurrentUser.GetID().Hex()},
				bson.M{"users": bson.M{"$elemMatch": bson.M{"$eq": c.CurrentUser.GetID().Hex()}}},
			}},
			bson.M{"is_visible": true},
			bson.M{"resource_type": bson.M{"$ne": core.ModelTypeChallenge}},
		}}
		if page <= 1 {
			page = 1
		}
		var pipe = mgomap.Aggregate{}.Match(match).Sort(bson.M{"updated_at": -1}).Skip((page - 1) * core.LimitRewards).Limit(core.LimitRewards)

		if err := Reward.Pipe(pipe, &rewards); err != nil {
			return c.ErrorResponse(c.Message("error.notFound", "Rewards"), "No rewards Found", 400)
		}
		return c.SuccessResponse(rewards, "success", core.ModelsType[core.ModelReward], serializers.RewardSerializer{Lang: c.Request.Locale})

	}
	return c.ServerErrorResponse()
}

/*
func ChangeRewardsModel() {

	var rewards []models.RewardBK
	if Reward, ok := app.Mapper.GetModel(&models.RewardBK{}); ok {
		if err := Reward.All().Exec(&rewards); err == nil {
			for _, v := range rewards {

				if nameStr, ok := v.Name.(string); ok {
					v.Name = map[string]interface{}{"default": nameStr}
				}
				if descrStr, ok := v.Description.(string); ok {
					v.Description = map[string]interface{}{"default": descrStr}
				}
				Reward.Update(&v)
			}
		}
	}
}
*/
