package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"spyc_backend/app/models"
	"strconv"
	"time"

	"github.com/revel/revel/testing"
)

// SpycAPIV2Test is the struct used for testing
type UserControllerTest struct {
	testing.TestSuite
}

// TestpersonalData is a representation of PersonalData struct for testing
type TestpersonalData struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Gender    string      `json:"gender"`
	BirthDate interface{} `json:"birth_date"`
	Address   interface{} `json:"address"`
}

// UserTest is a representation of User struct for teting
type UserTest struct {
	Email        string     `json:"email" validate:"nonzero"`
	UserName     string     `json:"user_name"`
	ImageURL     string     `json:"image_url"`
	Geolocation  models.Geo `json:"geolocation"`
	Token        string     `json:"-"`                  // Not saved field
	Password     string     `json:"password,omitempty"` // Not saved field
	PasswordHash string     `json:"-"`
	Received     time.Time  `json:"-"`
}

// UserTestResponse is the response the API returns
type UserTestResponse struct {
	UserTest UserTest `json:"user"`
}

// Before Run this before a request
func (t *UserControllerTest) Before() {
	println("Set up")
}

// TestPingApi tests PingApi from api/v2
func (t *UserControllerTest) TestPingApi() {
	var uri = "http://localhost:9000/api/v2/pingapi"
	customRequest := t.GetCustom(uri)
	customRequest.Send()

	log.Print(t.Response.StatusCode)

	// Validate response status to 200
	t.Assert(t.Response.StatusCode == 200)
}

// TestShowUser tests Show controller from api/v2
func (t *UserControllerTest) TestShowUser() {

	// Variables used for the post request
	var userIDTest = "5914a810e7f4f829771f141e"
	var uri = "http://localhost:9000/api/v2/users/" + userIDTest
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NzA2ODYsImlkIjoiNTkxMGIwMWVjNjNmMDEwYmM0MDQ5Mzk2IiwiQWN0aW9uIjoiYXV0aCJ9.l-IVJxM-t8_H0ROkadFOKKSOD0Fp7K_cM9OOGLOMEo8"

	// Generate custom request
	customRequest := t.GetCustom(uri)

	// Add header needed for the API
	customRequest.Header.Add("Authorization", "Bearer "+token)

	// Send Get Request
	customRequest.Send()

	// Parsing the response from server.
	var user UserTestResponse
	err := json.Unmarshal(t.ResponseBody, &user)
	if err != nil {
		t.Assert(false)
		log.Printf("Error parsing user ResponseBody")
	}

	// Register the time when new user data are received. We are not using "Expire"
	// parameter from persona test server because then we'll have to synchronise our clock.
	user.UserTest.Received = time.Now()
	log.Printf(user.UserTest.UserName)
}

// TestCreateUser tests Create method from UsersController
func (t *UserControllerTest) TestCreateUser() {
	var total = 1000

	// Variables used for the post request
	var user = UserTest{}
	var userResponse = UserTestResponse{}
	var uri = "/api/v2/users/"
	var contentType = "application/json"

	for i := 0; i < total; i++ {
		// fill userstruct
		user.Email = "test" + strconv.Itoa(i) + "@mail.com"
		user.UserName = "testUserName" + strconv.Itoa(i)
		user.Password = "password"
		user.Geolocation = models.Geo{
			Coordinates: []float64{0.0, 0.0},
		}

		// Create json
		userJSON, err := json.Marshal(&user)
		if err != nil {
			log.Print(err)
		}

		// Send json as new buffer to the given uri and
		t.Post(uri, contentType, bytes.NewBuffer(userJSON))

		// Print response
		log.Print(string(t.ResponseBody))

		err = json.Unmarshal(t.ResponseBody, &userResponse)
		if err != nil {
			t.Assert(false)
			log.Print("Error parsing user responseBody")
		}

		time.Sleep(100 * time.Millisecond)
	}

}
