// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tDBaseController struct {}
var DBaseController tDBaseController


func (_ tDBaseController) UnauthorizedResponse(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DBaseController.UnauthorizedResponse", args).URL
}

func (_ tDBaseController) ForbiddenResponse(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DBaseController.ForbiddenResponse", args).URL
}

func (_ tDBaseController) ServerErrorResponse(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DBaseController.ServerErrorResponse", args).URL
}

func (_ tDBaseController) FileResponse(
		file []byte,
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "file", file)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DBaseController.FileResponse", args).URL
}


type tBaseController struct {}
var BaseController tBaseController


func (_ tBaseController) UnauthorizedResponse(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("BaseController.UnauthorizedResponse", args).URL
}

func (_ tBaseController) ForbiddenResponse(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("BaseController.ForbiddenResponse", args).URL
}

func (_ tBaseController) ServerErrorResponse(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("BaseController.ServerErrorResponse", args).URL
}

func (_ tBaseController) FileResponse(
		file []byte,
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "file", file)
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("BaseController.FileResponse", args).URL
}


type tJobs struct {}
var Jobs tJobs


func (_ tJobs) Status(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Jobs.Status", args).URL
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).URL
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).URL
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).URL
}

func (_ tTestRunner) Suite(
		suite string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	return revel.MainRouter.Reverse("TestRunner.Suite", args).URL
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).URL
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).URL
}


type tDCampaignsController struct {}
var DCampaignsController tDCampaignsController


func (_ tDCampaignsController) Index(
		page int,
		quantity int,
		search string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "quantity", quantity)
	revel.Unbind(args, "search", search)
	return revel.MainRouter.Reverse("DCampaignsController.Index", args).URL
}

func (_ tDCampaignsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DCampaignsController.Show", args).URL
}

func (_ tDCampaignsController) Create(
		sponsorID string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sponsorID", sponsorID)
	return revel.MainRouter.Reverse("DCampaignsController.Create", args).URL
}

func (_ tDCampaignsController) ActivateCampaign(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DCampaignsController.ActivateCampaign", args).URL
}

func (_ tDCampaignsController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DCampaignsController.Delete", args).URL
}

func (_ tDCampaignsController) GetCampaignsBySponsor(
		sponsorID string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sponsorID", sponsorID)
	return revel.MainRouter.Reverse("DCampaignsController.GetCampaignsBySponsor", args).URL
}


type tDMissionsController struct {}
var DMissionsController tDMissionsController


func (_ tDMissionsController) Index(
		page int,
		quantity int,
		search string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "quantity", quantity)
	revel.Unbind(args, "search", search)
	return revel.MainRouter.Reverse("DMissionsController.Index", args).URL
}

func (_ tDMissionsController) New(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DMissionsController.New", args).URL
}

func (_ tDMissionsController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DMissionsController.Create", args).URL
}

func (_ tDMissionsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DMissionsController.Show", args).URL
}

func (_ tDMissionsController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DMissionsController.Update", args).URL
}

func (_ tDMissionsController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DMissionsController.Delete", args).URL
}

func (_ tDMissionsController) ActivateMission(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DMissionsController.ActivateMission", args).URL
}

func (_ tDMissionsController) GetMissionsByCampaign(
		idCampaign string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idCampaign", idCampaign)
	return revel.MainRouter.Reverse("DMissionsController.GetMissionsByCampaign", args).URL
}


type tDRewardsController struct {}
var DRewardsController tDRewardsController


func (_ tDRewardsController) Index(
		page int,
		quantity int,
		idCampaign string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "quantity", quantity)
	revel.Unbind(args, "idCampaign", idCampaign)
	return revel.MainRouter.Reverse("DRewardsController.Index", args).URL
}

func (_ tDRewardsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DRewardsController.Show", args).URL
}

func (_ tDRewardsController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DRewardsController.Create", args).URL
}

func (_ tDRewardsController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DRewardsController.Update", args).URL
}

func (_ tDRewardsController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DRewardsController.Delete", args).URL
}

func (_ tDRewardsController) ActivateReward(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DRewardsController.ActivateReward", args).URL
}


type tDAppPathsController struct {}
var DAppPathsController tDAppPathsController


func (_ tDAppPathsController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DAppPathsController.Index", args).URL
}

func (_ tDAppPathsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DAppPathsController.Show", args).URL
}

func (_ tDAppPathsController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DAppPathsController.Create", args).URL
}

func (_ tDAppPathsController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DAppPathsController.Update", args).URL
}

func (_ tDAppPathsController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DAppPathsController.Delete", args).URL
}


type tDGamesController struct {}
var DGamesController tDGamesController


func (_ tDGamesController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DGamesController.Index", args).URL
}

func (_ tDGamesController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DGamesController.Show", args).URL
}

func (_ tDGamesController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DGamesController.Create", args).URL
}

func (_ tDGamesController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DGamesController.Update", args).URL
}

func (_ tDGamesController) ActivateGame(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DGamesController.ActivateGame", args).URL
}

func (_ tDGamesController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DGamesController.Delete", args).URL
}


type tDSessionsController struct {}
var DSessionsController tDSessionsController


func (_ tDSessionsController) New(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DSessionsController.New", args).URL
}

func (_ tDSessionsController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DSessionsController.Create", args).URL
}

func (_ tDSessionsController) Delete(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DSessionsController.Delete", args).URL
}


type tDCashHuntController struct {}
var DCashHuntController tDCashHuntController


func (_ tDCashHuntController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DCashHuntController.Index", args).URL
}

func (_ tDCashHuntController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DCashHuntController.Create", args).URL
}

func (_ tDCashHuntController) Picture(
		idMission string,
		idTarget string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idMission", idMission)
	revel.Unbind(args, "idTarget", idTarget)
	return revel.MainRouter.Reverse("DCashHuntController.Picture", args).URL
}


type tDTargetsController struct {}
var DTargetsController tDTargetsController


func (_ tDTargetsController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DTargetsController.Index", args).URL
}

func (_ tDTargetsController) Create(
		missionID string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "missionID", missionID)
	return revel.MainRouter.Reverse("DTargetsController.Create", args).URL
}

func (_ tDTargetsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DTargetsController.Show", args).URL
}

func (_ tDTargetsController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DTargetsController.Update", args).URL
}

func (_ tDTargetsController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DTargetsController.Delete", args).URL
}

func (_ tDTargetsController) AddQuestion(
		idTarget string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idTarget", idTarget)
	return revel.MainRouter.Reverse("DTargetsController.AddQuestion", args).URL
}

func (_ tDTargetsController) GenerateQR(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DTargetsController.GenerateQR", args).URL
}

func (_ tDTargetsController) GetTargetsByMission(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DTargetsController.GetTargetsByMission", args).URL
}

func (_ tDTargetsController) Validation(
		page int,
		quantity int,
		targetType string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "quantity", quantity)
	revel.Unbind(args, "targetType", targetType)
	return revel.MainRouter.Reverse("DTargetsController.Validation", args).URL
}

func (_ tDTargetsController) Validate(
		idMission string,
		idSubscription string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idMission", idMission)
	revel.Unbind(args, "idSubscription", idSubscription)
	return revel.MainRouter.Reverse("DTargetsController.Validate", args).URL
}

func (_ tDTargetsController) GetSubscriptionsByTarget(
		idTarget string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idTarget", idTarget)
	return revel.MainRouter.Reverse("DTargetsController.GetSubscriptionsByTarget", args).URL
}

func (_ tDTargetsController) UpdateTargetsOrder(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DTargetsController.UpdateTargetsOrder", args).URL
}


type tDashboardController struct {}
var DashboardController tDashboardController


func (_ tDashboardController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DashboardController.Index", args).URL
}


type tDSponsorsController struct {}
var DSponsorsController tDSponsorsController


func (_ tDSponsorsController) Index(
		page int,
		quantity int,
		search string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	revel.Unbind(args, "quantity", quantity)
	revel.Unbind(args, "search", search)
	return revel.MainRouter.Reverse("DSponsorsController.Index", args).URL
}

func (_ tDSponsorsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DSponsorsController.Show", args).URL
}

func (_ tDSponsorsController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DSponsorsController.Create", args).URL
}

func (_ tDSponsorsController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DSponsorsController.Update", args).URL
}

func (_ tDSponsorsController) ActivateSponsor(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DSponsorsController.ActivateSponsor", args).URL
}

func (_ tDSponsorsController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DSponsorsController.Delete", args).URL
}


type tDUsersController struct {}
var DUsersController tDUsersController


func (_ tDUsersController) Index(
		page int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("DUsersController.Index", args).URL
}

func (_ tDUsersController) Create(
		role string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "role", role)
	return revel.MainRouter.Reverse("DUsersController.Create", args).URL
}

func (_ tDUsersController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("DUsersController.Show", args).URL
}

func (_ tDUsersController) ResetPassword(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DUsersController.ResetPassword", args).URL
}

func (_ tDUsersController) ResetPasswordSuccess(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DUsersController.ResetPasswordSuccess", args).URL
}

func (_ tDUsersController) ChangePassword(
		password string,
		confirmation string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "password", password)
	revel.Unbind(args, "confirmation", confirmation)
	return revel.MainRouter.Reverse("DUsersController.ChangePassword", args).URL
}


type tDSChallengesController struct {}
var DSChallengesController tDSChallengesController


func (_ tDSChallengesController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("DSChallengesController.Index", args).URL
}

func (_ tDSChallengesController) GetStats(
		status string,
		country string,
		statstype string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "status", status)
	revel.Unbind(args, "country", country)
	revel.Unbind(args, "statstype", statstype)
	return revel.MainRouter.Reverse("DSChallengesController.GetStats", args).URL
}


type tSaasController struct {}
var SaasController tSaasController


func (_ tSaasController) Auth(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SaasController.Auth", args).URL
}

func (_ tSaasController) User(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SaasController.User", args).URL
}

func (_ tSaasController) Archive(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SaasController.Archive", args).URL
}

func (_ tSaasController) GetRoster(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SaasController.GetRoster", args).URL
}

func (_ tSaasController) NewRoster(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SaasController.NewRoster", args).URL
}

func (_ tSaasController) UpdateRoster(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SaasController.UpdateRoster", args).URL
}


type tCoinsController struct {}
var CoinsController tCoinsController


func (_ tCoinsController) SendTransaction(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("CoinsController.SendTransaction", args).URL
}


type tWebGamesController struct {}
var WebGamesController tWebGamesController


func (_ tWebGamesController) StartWebGame(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("WebGamesController.StartWebGame", args).URL
}

func (_ tWebGamesController) ValidateWebGame(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("WebGamesController.ValidateWebGame", args).URL
}

func (_ tWebGamesController) UpdateWebGameScore(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("WebGamesController.UpdateWebGameScore", args).URL
}

func (_ tWebGamesController) GetWebGameStatus(
		idTarget string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idTarget", idTarget)
	return revel.MainRouter.Reverse("WebGamesController.GetWebGameStatus", args).URL
}

func (_ tWebGamesController) GetWebGames(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("WebGamesController.GetWebGames", args).URL
}


type tChallengesController struct {}
var ChallengesController tChallengesController


func (_ tChallengesController) CreateChallenge(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ChallengesController.CreateChallenge", args).URL
}

func (_ tChallengesController) SaveChallengeFile(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChallengesController.SaveChallengeFile", args).URL
}

func (_ tChallengesController) GetChallengeInfo(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChallengesController.GetChallengeInfo", args).URL
}

func (_ tChallengesController) StartChallenge(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ChallengesController.StartChallenge", args).URL
}

func (_ tChallengesController) ValidateChallenge(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ChallengesController.ValidateChallenge", args).URL
}

func (_ tChallengesController) UpdateChallengeScore(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ChallengesController.UpdateChallengeScore", args).URL
}

func (_ tChallengesController) GetChallengeUserStatus(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChallengesController.GetChallengeUserStatus", args).URL
}

func (_ tChallengesController) GetChallengePlayerStatus(
		id string,
		idPlayer string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "idPlayer", idPlayer)
	return revel.MainRouter.Reverse("ChallengesController.GetChallengePlayerStatus", args).URL
}

func (_ tChallengesController) GetChallengeStatus(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChallengesController.GetChallengeStatus", args).URL
}

func (_ tChallengesController) NotifyChallengePlayers(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChallengesController.NotifyChallengePlayers", args).URL
}

func (_ tChallengesController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChallengesController.Delete", args).URL
}


type tAdvertisementController struct {}
var AdvertisementController tAdvertisementController


func (_ tAdvertisementController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("AdvertisementController.Create", args).URL
}

func (_ tAdvertisementController) UploadAttachment(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("AdvertisementController.UploadAttachment", args).URL
}

func (_ tAdvertisementController) GetAdvertisementsBySponsor(
		idSponsor string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idSponsor", idSponsor)
	return revel.MainRouter.Reverse("AdvertisementController.GetAdvertisementsBySponsor", args).URL
}


type tNotificationsController struct {}
var NotificationsController tNotificationsController


func (_ tNotificationsController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("NotificationsController.Index", args).URL
}

func (_ tNotificationsController) MarkAsReceived(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("NotificationsController.MarkAsReceived", args).URL
}

func (_ tNotificationsController) MarkAsSeen(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("NotificationsController.MarkAsSeen", args).URL
}

func (_ tNotificationsController) MarkAsCompleted(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("NotificationsController.MarkAsCompleted", args).URL
}


type tHistoriesController struct {}
var HistoriesController tHistoriesController


func (_ tHistoriesController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("HistoriesController.Index", args).URL
}

func (_ tHistoriesController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("HistoriesController.Create", args).URL
}

func (_ tHistoriesController) AddView(
		viewedtime int,
		id string,
		completed string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "viewedtime", viewedtime)
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "completed", completed)
	return revel.MainRouter.Reverse("HistoriesController.AddView", args).URL
}


type tInvitesController struct {}
var InvitesController tInvitesController


func (_ tInvitesController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("InvitesController.Index", args).URL
}

func (_ tInvitesController) Respond(
		id string,
		respond string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "respond", respond)
	return revel.MainRouter.Reverse("InvitesController.Respond", args).URL
}


type tPINController struct {}
var PINController tPINController


func (_ tPINController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("PINController.Index", args).URL
}

func (_ tPINController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("PINController.Create", args).URL
}

func (_ tPINController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("PINController.Update", args).URL
}

func (_ tPINController) Delete(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("PINController.Delete", args).URL
}

func (_ tPINController) PinClicked(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("PINController.PinClicked", args).URL
}

func (_ tPINController) PinViewed(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("PINController.PinViewed", args).URL
}


type tTargetsController struct {}
var TargetsController tTargetsController


func (_ tTargetsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("TargetsController.Show", args).URL
}

func (_ tTargetsController) GetCurrentTargetByMissionID(
		idMission string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "idMission", idMission)
	return revel.MainRouter.Reverse("TargetsController.GetCurrentTargetByMissionID", args).URL
}

func (_ tTargetsController) SubscribeToTarget(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("TargetsController.SubscribeToTarget", args).URL
}

func (_ tTargetsController) UnsubscribeFromTarget(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("TargetsController.UnsubscribeFromTarget", args).URL
}

func (_ tTargetsController) ValidateReward(
		file []byte,
		idMission string,
		idTarget string,
		data string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "file", file)
	revel.Unbind(args, "idMission", idMission)
	revel.Unbind(args, "idTarget", idTarget)
	revel.Unbind(args, "data", data)
	return revel.MainRouter.Reverse("TargetsController.ValidateReward", args).URL
}


type tSanctionsController struct {}
var SanctionsController tSanctionsController


func (_ tSanctionsController) SactionUser(
		sanctionType string,
		userID string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "sanctionType", sanctionType)
	revel.Unbind(args, "userID", userID)
	return revel.MainRouter.Reverse("SanctionsController.SactionUser", args).URL
}


type tColorsController struct {}
var ColorsController tColorsController


func (_ tColorsController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ColorsController.Index", args).URL
}


type tSMSController struct {}
var SMSController tSMSController


func (_ tSMSController) RequestNew(
		phone string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "phone", phone)
	return revel.MainRouter.Reverse("SMSController.RequestNew", args).URL
}

func (_ tSMSController) VerifySMS(
		code string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "code", code)
	return revel.MainRouter.Reverse("SMSController.VerifySMS", args).URL
}


type tSessionsController struct {}
var SessionsController tSessionsController


func (_ tSessionsController) Create(
		account string,
		password string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "account", account)
	revel.Unbind(args, "password", password)
	return revel.MainRouter.Reverse("SessionsController.Create", args).URL
}


type tRewardsController struct {}
var RewardsController tRewardsController


func (_ tRewardsController) CollectReward(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("RewardsController.CollectReward", args).URL
}

func (_ tRewardsController) GetRewards(
		page int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("RewardsController.GetRewards", args).URL
}


type tSpyNewsController struct {}
var SpyNewsController tSpyNewsController


func (_ tSpyNewsController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SpyNewsController.Index", args).URL
}

func (_ tSpyNewsController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SpyNewsController.Create", args).URL
}

func (_ tSpyNewsController) Comment(
		id string,
		message string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "message", message)
	return revel.MainRouter.Reverse("SpyNewsController.Comment", args).URL
}

func (_ tSpyNewsController) React(
		id string,
		react string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "react", react)
	return revel.MainRouter.Reverse("SpyNewsController.React", args).URL
}


type tMissionsController struct {}
var MissionsController tMissionsController


func (_ tMissionsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("MissionsController.Show", args).URL
}

func (_ tMissionsController) ShowAllMissions(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("MissionsController.ShowAllMissions", args).URL
}

func (_ tMissionsController) GetSubscribedMissions(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("MissionsController.GetSubscribedMissions", args).URL
}

func (_ tMissionsController) ShowMissionsByLocation(
		lat float64,
		lng float64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "lat", lat)
	revel.Unbind(args, "lng", lng)
	return revel.MainRouter.Reverse("MissionsController.ShowMissionsByLocation", args).URL
}

func (_ tMissionsController) SubscribeToMission(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("MissionsController.SubscribeToMission", args).URL
}

func (_ tMissionsController) UnsubscribeFromMission(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("MissionsController.UnsubscribeFromMission", args).URL
}

func (_ tMissionsController) AddUserView(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("MissionsController.AddUserView", args).URL
}


type tMapsController struct {}
var MapsController tMapsController


func (_ tMapsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("MapsController.Show", args).URL
}

func (_ tMapsController) Create(
		coords string,
		name string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "coords", coords)
	revel.Unbind(args, "name", name)
	return revel.MainRouter.Reverse("MapsController.Create", args).URL
}


type tV2SaasController struct {}
var V2SaasController tV2SaasController


func (_ tV2SaasController) Info(
		jid string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "jid", jid)
	return revel.MainRouter.Reverse("V2SaasController.Info", args).URL
}


type tUsersController struct {}
var UsersController tUsersController


func (_ tUsersController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("UsersController.Index", args).URL
}

func (_ tUsersController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("UsersController.Show", args).URL
}

func (_ tUsersController) Create(
		body string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "body", body)
	return revel.MainRouter.Reverse("UsersController.Create", args).URL
}

func (_ tUsersController) Update(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("UsersController.Update", args).URL
}

func (_ tUsersController) Delete(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("UsersController.Delete", args).URL
}

func (_ tUsersController) ProfilePicture(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("UsersController.ProfilePicture", args).URL
}

func (_ tUsersController) ChangeUsername(
		username string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	return revel.MainRouter.Reverse("UsersController.ChangeUsername", args).URL
}

func (_ tUsersController) Language(
		lang string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "lang", lang)
	return revel.MainRouter.Reverse("UsersController.Language", args).URL
}

func (_ tUsersController) Coords(
		lat float64,
		lng float64,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "lat", lat)
	revel.Unbind(args, "lng", lng)
	return revel.MainRouter.Reverse("UsersController.Coords", args).URL
}

func (_ tUsersController) MessagingToken(
		token string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "token", token)
	return revel.MainRouter.Reverse("UsersController.MessagingToken", args).URL
}

func (_ tUsersController) VerifyUserName(
		account string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "account", account)
	return revel.MainRouter.Reverse("UsersController.VerifyUserName", args).URL
}

func (_ tUsersController) FriendRequest(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("UsersController.FriendRequest", args).URL
}

func (_ tUsersController) Friends(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("UsersController.Friends", args).URL
}

func (_ tUsersController) ResetPassword(
		email string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	return revel.MainRouter.Reverse("UsersController.ResetPassword", args).URL
}


type tChestsController struct {}
var ChestsController tChestsController


func (_ tChestsController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChestsController.Show", args).URL
}

func (_ tChestsController) Search(
		tags string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "tags", tags)
	return revel.MainRouter.Reverse("ChestsController.Search", args).URL
}

func (_ tChestsController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("ChestsController.Create", args).URL
}

func (_ tChestsController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChestsController.Update", args).URL
}

func (_ tChestsController) AddFile(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChestsController.AddFile", args).URL
}

func (_ tChestsController) RemoveFile(
		id string,
		file string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "file", file)
	return revel.MainRouter.Reverse("ChestsController.RemoveFile", args).URL
}

func (_ tChestsController) Lock(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("ChestsController.Lock", args).URL
}

func (_ tChestsController) Unlock(
		id string,
		pin string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "pin", pin)
	return revel.MainRouter.Reverse("ChestsController.Unlock", args).URL
}


type tImagesController struct {}
var ImagesController tImagesController


func (_ tImagesController) GetResource(
		gameType string,
		size string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "gameType", gameType)
	revel.Unbind(args, "size", size)
	return revel.MainRouter.Reverse("ImagesController.GetResource", args).URL
}


type tStatsController struct {}
var StatsController tStatsController


func (_ tStatsController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("StatsController.Index", args).URL
}

func (_ tStatsController) Replicate(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("StatsController.Replicate", args).URL
}


type tBlacklistController struct {}
var BlacklistController tBlacklistController


func (_ tBlacklistController) Blacklist(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("BlacklistController.Blacklist", args).URL
}

func (_ tBlacklistController) AddToBlackList(
		jid string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "jid", jid)
	return revel.MainRouter.Reverse("BlacklistController.AddToBlackList", args).URL
}

func (_ tBlacklistController) RemoveFromBlacklist(
		jid string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "jid", jid)
	return revel.MainRouter.Reverse("BlacklistController.RemoveFromBlacklist", args).URL
}

func (_ tBlacklistController) CheckJIDInBlacklist(
		jid string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "jid", jid)
	return revel.MainRouter.Reverse("BlacklistController.CheckJIDInBlacklist", args).URL
}


type tSearchController struct {}
var SearchController tSearchController


func (_ tSearchController) TimeFormat(
		format string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "format", format)
	return revel.MainRouter.Reverse("SearchController.TimeFormat", args).URL
}

func (_ tSearchController) GetPeople(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SearchController.GetPeople", args).URL
}

func (_ tSearchController) Search(
		account string,
		saas bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "account", account)
	revel.Unbind(args, "saas", saas)
	return revel.MainRouter.Reverse("SearchController.Search", args).URL
}

func (_ tSearchController) InvitePeople(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("SearchController.InvitePeople", args).URL
}


type tEjabberdController struct {}
var EjabberdController tEjabberdController


func (_ tEjabberdController) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("EjabberdController.Index", args).URL
}

func (_ tEjabberdController) Show(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("EjabberdController.Show", args).URL
}

func (_ tEjabberdController) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("EjabberdController.Create", args).URL
}

func (_ tEjabberdController) Update(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("EjabberdController.Update", args).URL
}

func (_ tEjabberdController) AddMember(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("EjabberdController.AddMember", args).URL
}

func (_ tEjabberdController) UpdateMember(
		id string,
		jid string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "jid", jid)
	return revel.MainRouter.Reverse("EjabberdController.UpdateMember", args).URL
}

func (_ tEjabberdController) RemoveMember(
		id string,
		jid string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "jid", jid)
	return revel.MainRouter.Reverse("EjabberdController.RemoveMember", args).URL
}


type tCommentsController struct {}
var CommentsController tCommentsController


func (_ tCommentsController) React(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("CommentsController.React", args).URL
}


type tSectionsController struct {}
var SectionsController tSectionsController


func (_ tSectionsController) Create(
		t int,
		name string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "t", t)
	revel.Unbind(args, "name", name)
	return revel.MainRouter.Reverse("SectionsController.Create", args).URL
}

func (_ tSectionsController) Replicate(
		id string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("SectionsController.Replicate", args).URL
}

func (_ tSectionsController) Info(
		section string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "section", section)
	return revel.MainRouter.Reverse("SectionsController.Info", args).URL
}

func (_ tSectionsController) NotifyAdvertisingSeen(
		section string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "section", section)
	return revel.MainRouter.Reverse("SectionsController.NotifyAdvertisingSeen", args).URL
}


type tMessagesController struct {}
var MessagesController tMessagesController


func (_ tMessagesController) Index(
		user string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "user", user)
	return revel.MainRouter.Reverse("MessagesController.Index", args).URL
}

func (_ tMessagesController) UploadFile(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("MessagesController.UploadFile", args).URL
}

func (_ tMessagesController) DeleteFile(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("MessagesController.DeleteFile", args).URL
}


