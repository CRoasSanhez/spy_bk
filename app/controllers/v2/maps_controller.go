package v2

import (
	"io/ioutil"
	"net/http"
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"log"

	"github.com/revel/revel"
)

// MapsController controller
type MapsController struct {
	BaseController
}

// Show get a location image
func (c *MapsController) Show(id string) revel.Result {
	var myMap models.Map

	if err := models.GetDocument(id, &myMap); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	res, err := http.Get(myMap.URL)
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	c.Response.Out.Header().Set("map_id", revel.HTTPAddr+"/v2/maps/"+myMap.ID.Hex())
	return c.FileResponse(body, myMap.ID.Hex())
}

// Create make a mew map
func (c *MapsController) Create(coords, name string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(coords) <= 0 {
		return c.ErrorResponse(nil, "Invalid Coords", 400)
	}

	var myMap models.Map
	var baseURL = "https://maps.googleapis.com/maps/api/staticmap?center="
	var fullURL = baseURL + coords + "&zoom=16&size=400x400&markers=color:blue%7Clabel:%7C" + coords + "&key=" + core.GoogleStaticMapsAPIKey

	myMap.Name = name
	myMap.Coords = coords
	myMap.URL = fullURL

	if err := models.CreateDocument(&myMap); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	res, err := http.Get(fullURL)
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return c.ErrorResponse(err, err.Error(), 400)
	}

	c.Response.Out.Header().Set("map_id", revel.HTTPAddr+"/v2/maps/"+myMap.ID.Hex())
	return c.FileResponse(body, myMap.ID.Hex())
}
