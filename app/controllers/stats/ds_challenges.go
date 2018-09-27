package stats

import (
	"spyc_backend/app"
	dashboard "spyc_backend/app/controllers/dashboard"
	"spyc_backend/app/models"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

type DSChallengesController struct {
	dashboard.DBaseController
}

// Index returns the stats dashboard
func (c DSChallengesController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}
	return c.Render()
}

// GetStats [/spyc_stats/f/:statstype] GET
// returns a list of challenges based on the given filter
func (c DSChallengesController) GetStats(status, country, statstype string) revel.Result {

	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var pipe = mgomap.Aggregate{}
	var stats []interface{}
	var group bson.M

	if status != "" {
		pipe = pipe.Match(bson.M{"status.name": bson.M{"$in": []string{status}}})
	}

	pipe = pipe.LookUp("users", "user_id", "_id", "user").Add(bson.M{"$unwind": "$user"}).Add(bson.M{"$unwind": "$players"})

	if country != "" {
		pipe.Match(bson.M{"user.personal_data.address.country": country})
	}

	if status != "" {
		group = bson.M{"$group": bson.M{
			"_id": bson.M{
				//"c_type":     "$type",
				"c_type":     "$status.name",
				"country":    "$user.personal_data.address.country",
				"created_at": bson.M{"$dateToString": bson.M{"format": "%Y-%m-15", "date": "$created_at"}},
			},
			"count": bson.M{"$sum": 1},
		}}
	} else {
		group = bson.M{"$group": bson.M{
			"_id": bson.M{
				"c_type":     "$type",
				"country":    "$user.personal_data.address.country",
				"created_at": bson.M{"$dateToString": bson.M{"format": "%Y-%m-15", "date": "$created_at"}},
			},
			"count": bson.M{"$sum": 1},
		}}
	}

	var project = bson.M{"$project": bson.M{
		"_id": 0, "c_type": "$_id.c_type", "country": "$_id.country", "created_at": "$_id.created_at", "count": "$count",
	}}

	pipe = pipe.Add(group).Add(project).Sort(bson.M{"c_type": 1, "created_at": 1})

	if Challenge, ok := app.Mapper.GetModel(&models.Challenge{}); ok {

		if err := Challenge.Pipe(pipe, &stats); err != nil {
			revel.ERROR.Printf("ERROR FIND Stats --- %s", err.Error())
			return c.ErrorResponse(c.Message("error.notFound", stats, ""), "No challenges Found", 400)
		}
		return c.SuccessResponse(stats, "success", 200, nil)
	}
	return c.ServerErrorResponse()
}
