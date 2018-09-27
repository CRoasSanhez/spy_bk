package app

import (
	"strconv"
	"time"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// Mapper storage mgomap object, this is responsible for handling request for MongoDB
var Mapper *mgomap.Mapper

// InitDB initialize DB
func InitDB() {
	if !revel.DevMode {
		Mapper = &mgomap.Mapper{
			Hosts: []string{
				"spychatterprod-shard-00-00-b0z3q.mongodb.net:27017",
				"spychatterprod-shard-00-01-b0z3q.mongodb.net:27017",
				"spychatterprod-shard-00-02-b0z3q.mongodb.net:27017",
			},
			DatabaseName:    revel.Config.StringDefault("database.name", ""),
			DatabaseUser:    "spyc_back",
			DatabasePass:    "7NO6zs6myp6cAbnT",
			DatabaseReplica: "SpychatterProd-shard-0",
			Atlas:           true,
			TimeOut:         5 * time.Second,
		}
	} else {
		Mapper = &mgomap.Mapper{
			Hosts:        []string{revel.Config.StringDefault("database.host", "127.0.0.1")},
			DatabaseName: revel.Config.StringDefault("database.name", ""),
			TimeOut:      5 * time.Second,
		}
	}

	Mapper.Connect()
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		CORSFilter,                    // CORS
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.TemplateFuncs["N"] = func(size int) []string {
		return make([]string, size)
	}

	revel.TemplateFuncs["add"] = func(a, b int) int {
		return a + b
	}

	revel.TemplateFuncs["to_s"] = func(a int) string {
		return strconv.Itoa(a)
	}

	// Agregamos el formato de fecha
	revel.TimeFormats = append(revel.TimeFormats, "02/01/2006 15:04")

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter ...
// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

// CORSFilter ...
var CORSFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Out.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	c.Response.Out.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Stop here for a Preflighted OPTIONS request.
	if c.Request.Method == "OPTIONS" {
		return
	}

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
