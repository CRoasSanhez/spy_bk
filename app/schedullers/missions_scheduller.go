package schedullers

import (
	"spyc_backend/app"
	controllers "spyc_backend/app/controllers/dashboard"
	"spyc_backend/app/models"
	"strconv"
	"time"

	"github.com/Reti500/mgomap"

	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app/core"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

type ChangeMissionsStatus struct{}
type SpyScheduller struct{}

// Run activates the missions in stats INIT
func (jms ChangeMissionsStatus) Run() {

	var dbController = controllers.DBaseController{
		Controller:  revel.NewControllerEmpty(),
		CurrentUser: models.User{},
	}

	var missionController = controllers.DMissionsController{dbController}

	if Mission, ok := app.Mapper.GetModel(&models.Mission{}); ok {
		var missions []models.Mission
		var now = core.ChangeLocalTimeToUTCZone(time.Now())

		var match = bson.M{"$match": bson.M{
			"status.name": bson.M{"$in": []string{core.StatusInit, core.StatusActive}},
		}}

		var proj = bson.M{
			"$project": bson.M{
				"_id":    1,
				"status": 1, "title": 1, "description": 1, "type": 1, "start_date": 1, "end_date": 1, "countries": 1, "langs": 1, "geolocation": 1,

				"date_diff": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []string{"$status.name", "init"}},
						"then": bson.M{"$subtract": []interface{}{now, "$start_date"}},
						"else": bson.M{"$subtract": []interface{}{"$end_date", now}},
					},
				},
			},
		}

		var pipe = mgomap.Aggregate{}.Add(match).Add(proj)

		if err := Mission.Pipe(pipe, &missions); err != nil {
			revel.ERROR.Printf("ERROR ACTIVATE JOB Missions --- %s", err.Error())
			return
		}

		for i := 0; i < len(missions); i++ {

			// If missions has to be deactivated
			if missions[i].Status.Name == core.StatusActive {
				if missions[i].DateDiff >= 0 && missions[i].DateDiff <= core.SchedullerDelayMiliseconds {
					if missions[i].Complete() {
						dbController.NotifyMissionComplete(missions[i])
					}
				}
			} else {
				// If mission has to be activated
				// If difference is bigger then 0 then mission has a delay with its deactivation
				if missions[i].DateDiff >= 0 && missions[i].DateDiff <= core.SchedullerDelayMiliseconds {
					missionController.MissionActivation(missions[i])
				}
			}

		}
	}
}

// Run does what it has to do
func (jul SpyScheduller) Run() {

	var user models.User
	var err error

	if err = models.GetDocumentBy(bson.M{"email": "delgadojohanna@yahoo.com"}, &user); err != nil {
		revel.ERROR.Printf("ERROR Finding user --- %s", err.Error())
		return
	}
	var history = &struct {
		Coordinates []float64 `bson:"coordinates"`
		Date        time.Time `bson:"date"`
	}{
		Coordinates: user.Geolocation.Coordinates,
		Date:        time.Now(),
	}

	if err = models.UpdateByQuery(&user,
		bson.M{"_id": user.GetID()},
		bson.M{"$push": bson.M{"location_history": history}}, false); err != nil {
		revel.ERROR.Printf("ERROR UPDATE User Location %s --- %s", user.GetID(), err.Error())
	}

}

// init initializes the scheduller
func init() {
	revel.OnAppStart(func() {

		if revel.RunMode == "prod" {
			// Activate Mission from status INACTIVE
			// run in the 0 min to 59min and SchedullerTimeMinutes thereafter from 0 to 23hrs
			jobs.Schedule("0 0-59/"+strconv.Itoa(core.SchedullerTimeMinutes)+" 0-23 * * *", ChangeMissionsStatus{})

			// run in the 0 min to 59min and 30min thereafter from 0 to 23hrs
			jobs.Schedule("0 0-59/30 0-23 * * *", SpyScheduller{})
		} else {
			// run in the 0 min to 59min and SchedullerTimeMinutes thereafter from 0 to 23hrs
			jobs.Schedule("0 0-59/"+strconv.Itoa(core.SchedullerTimeMinutes)+" 0-23 * * *", ChangeMissionsStatus{})
		}
	})
}
