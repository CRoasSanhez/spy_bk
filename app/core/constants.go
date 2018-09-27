package core

import "github.com/revel/revel"

// App constants
const (
	AccessToken     = "access_token"
	ActionAuth      = "auth"
	ActionClick     = "click"
	ActionEjabberd  = "ejabberd"
	ActionDashboard = "dashboard"
	ActionManage    = "manage"
	ActionResetPass = "reset_pass"
	ActionStart     = "start"
	ActionLock      = "locked"
	ActionUnlock    = "unlock"
	ActionUpdate    = "update"
	ActionView      = "view"
	ApplicationJSON = "application/json"

	// Header
	Bearer       = "Bearer "
	UserLocation = "Uloc "
	ChestAccess  = "ChestAccess"
	Compression  = "Compression"

	// Defaults
	DefaultLimitDocuments = 20
	DefaultSigningTime    = 1440 // time in minutes ~ 1 day

	// Games type for both games and Targets
	GameTypeCheckIn = "checkin"
	GameTypeGame    = "webgame"
	GameTypeNIP     = "nip"
	GameTypeOptions = "options"
	GameTypePhoto   = "photo"
	GameTypeQR      = "qr"
	GameTypeText    = "text"

	// Constant Values
	BlurChallenge                  = 9
	BlurZero                       = 0.1
	LimitRewards                   = 10
	MaxCashHuntGamesIntent         = 15
	MaxDistanceForPin              = 1500
	MaxDistanceForAnswerValidation = 15
	MaxGameIntents                 = 2
	MaxPushRegIds                  = 990
	MaxRewardViews                 = 2
	MaxSubscriptions               = 1
	MaxTransactionAmount           = 500
	SchedullerTimeMinutes          = 5
	SchedullerDelayMiliseconds     = (SchedullerTimeMinutes * 60 * 1000) - 1

	// Mission Types
	TypeBroadcast            = "broadcast"
	TypePhysical             = "physical"
	TypeVirtual              = "virtual"
	TypeVirtualInternational = "virtual_international"
	TypeCountry              = "country"

	// Models Type
	ModelCashHunt       = "cash_hunt_response"
	ModelTypeChallenge  = "challenge"
	ModelFile           = "file"
	ModelGame           = "game"
	ModelGames          = "games"
	ModelPlayer         = "player"
	ModelQuestion       = "question"
	ModelReward         = "reward"
	ModelSimpleResponse = "simple_response"
	ModelStep           = "step"
	ModelWebGame        = "web_game"
	ModelWebGameInfo    = "web_game_info"
	ModelWebGameStart   = "web_game_start"

	// Notification Type
	NotificationFriendInvite      = "friend_invite"
	NotificationCashHuntCompleted = "cash_hunt_completed"
	NotificationCashHuntInvite    = "cash_hunt_invite"
	NotificationRewardObtained    = "reward_obtained"

	// PUSH Notification Types
	CashHuntNew           = "cash_hunt_new"
	CashHuntFinished      = "cash_hunt_finished"
	CashHuntReward        = "cash_hunt_reward"
	ChallengeReward       = "challenge_reward"
	FriendRequest         = "friend_request"
	FrienRequestConfirmed = "friend_request_confirmed"
	Message               = "chat_message"

	// Sanction Type
	SanctionScreenShoot = "scrsht"
	SanctionMissionQuit = "mssnqt"

	// Sections for permissions
	SectionDashboard = "dashboard"

	// Status
	StatusActive                  = "active"
	StatusBlocked                 = "blocked"
	StatusCreated                 = "created"
	StatusCompleted               = "completed"
	StatusCompletedNotWinner      = "completed_not_winner"
	StatusCompletedWinner         = "completed_winner"
	StatusWinnerMaxViewsReached   = "winner_max_views_reached"
	StatusCompletedWrong          = "completed_wrong"
	StatusEmailPending            = "pending_email"
	StatusExpired                 = "expired"
	StatusInactive                = "inactive"
	StatusIncomplete              = "incomplete"
	StatusInit                    = "init"
	StatusInvalidID               = "invalid_id"
	StatusLimitedTime             = "limited_time"
	StatusMaxIntentsReached       = "max_intents_reached"
	StatusMaxSubscriptionsReached = "max_subscriptions_reached"
	StatusNotFound                = "not_found"
	StatusNotSubscribed           = "not_subscribed"
	StatusObtained                = "reward_collected"
	StatusPendingCollection       = "pending_collection"
	StatusPendingConfirmation     = "pending_confirmation"
	StatusPendingValidation       = "pending_validation"
	StatusReceived                = "received"
	StatusSeen                    = "seen"
	StatusSMSSent                 = "sms_sent"
	StatusSMSConfirmed            = "sms_confirmed"
	StatusSuspended               = "suspended"
	StatusUnsubscribed            = "unsubscribed"
	StatusWrong                   = "wrong"

	// Response Status
	StatusError   = "error"
	StatusSuccess = "success"
	StatusUpdated = "updated"

	// App Sections
	SectionCashHunt = "cashhunt"

	// Formats
	MXTimeFormat   = "02/01/2006 15:04"
	MXTimeFormatTZ = "02/01/2006 15:04 -0700"
	JSONFormat     = "json"

	// Validations Models
	ValidationInvalidID = "invalid_id"
	ValidationMinLength = "min_length"
	ValidationMaxLength = "max_length"
	ValidationLength    = "length"
	ValidationMax       = "max"
	ValidationMin       = "min"
	ValidationRequired  = "required"

	// Wins Types
	WinTypeCashHunt  = "cash_hunt"
	WinTypeChatGames = "chat_games"

	// TrueString for true string parameter
	TrueString = "true"

	GoogleStaticMapsAPIKey = "AIzaSyCwgJvxL9s6yWxAC_n2_kV7pWq-O3YnhSY"
	GoogleMapsJSAPIKey     = "AIzaSyDAJBY8WwnriS81iciLLgWoTyV-uvzCdkc"

	// BasePath for webgames
	GameBasePathProd = "https://assets.spychatter.net/games/"
	GameBasePathDev  = "http://35.165.190.214:90/"

	// BasePath for api
	APIBasePathProd = "https://api.spychatter.net/"
	APIBasePathDev  = "http://35.165.190.214/"
	// APIBasePathDev = "http://localhost:9000/"
)

// Bytes contants
const (
	_        = iota // ignore first value by assigning to blank identifier
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
)

// Vars as contants
var (
	AccountStatus = map[string]int{
		"init":                 9001,
		"pending_confirmation": 9002,
		"pending_email":        9003,
		"active":               9005,
		"created":              9010,
		"sms_sent":             9006,
		"sms_confirmed":        9007,
		"blocked":              9050,
		"suspended":            9999,
	}

	// Status ued for both GAMES and STEPS
	GameStatus = map[string]int{
		StatusInit:              6001,
		StatusActive:            6002,
		StatusInactive:          6003,
		StatusLimitedTime:       6005,
		StatusCompleted:         6006,
		StatusIncomplete:        6007,
		StatusMaxIntentsReached: 6008,
		StatusUnsubscribed:      6009,
	}

	MissionTypes = map[string]string{
		TypeVirtual:              "Virtual",
		TypePhysical:             "Physical",
		TypeCountry:              "For Country",
		TypeVirtualInternational: "Virtual Internacional",
		TypeBroadcast:            "Broadcast",
	}

	ModelStatus = map[string]int{
		StatusInvalidID:     1000,
		ValidationMinLength: 1001,
		ValidationMaxLength: 1002,
		ValidationLength:    1003,
		ValidationMax:       1004,
		ValidationMin:       1005,
		ValidationRequired:  1006,
	}

	// Codes that define a specific model or structure
	ModelsType = map[string]int{
		ModelCashHunt:       2001,
		ModelStep:           2002,
		ModelGame:           2003,
		ModelReward:         2004,
		ModelQuestion:       2005,
		ModelSimpleResponse: 2006,
		ModelWebGameInfo:    2007,
		ModelGames:          2008,
		ModelTypeChallenge:  2009,
		ModelWebGameStart:   2010,
		ModelFile:           2011,
		ModelWebGame:        2012,
		ModelPlayer:         2013,
	}

	NotificationType = map[string]string{}

	SanctionType = map[string]int{
		SanctionScreenShoot: 7300,
		SanctionMissionQuit: 7301,
	}

	// Status used for both GAMES and STEPS subscription
	SubscriptionStatus = map[string]int{
		StatusInit:                  7000,
		StatusActive:                7001,
		StatusInactive:              7002,
		StatusCompleted:             7003,
		StatusCompletedNotWinner:    7004,
		StatusCompletedWinner:       7005,
		StatusIncomplete:            7006,
		StatusCompletedWrong:        7007,
		StatusWinnerMaxViewsReached: 7008,
		StatusPendingValidation:     7009,
	}

	ValidationStatus = map[string]int{
		StatusNotFound:                8001,
		StatusExpired:                 8002,
		StatusWrong:                   8003,
		StatusSuccess:                 8004,
		StatusCompleted:               8005,
		StatusMaxIntentsReached:       8006,
		StatusPendingCollection:       8007,
		StatusObtained:                8008,
		StatusMaxSubscriptionsReached: 8009,
		StatusError:                   8010,
		StatusNotSubscribed:           8011,
		StatusPendingValidation:       8012,
	}

	CurrentLocales = []string{
		"es",
		"en",
	}
)

// GetGameBasePath returns the game path based on the enviroment
func GetGameBasePath() string {
	if revel.RunMode == "prod" {
		return GameBasePathProd
	}
	return GameBasePathDev
}

// GetBasePath returns api base path
func GetBasePath() string {
	if revel.RunMode == "prod" {
		return APIBasePathProd
	}

	return APIBasePathDev
}

// GetDashboardPath returns base path for dashboard
func GetDashboardPath() string {
	return GetBasePath() + "spyc_admin/"
}
