package dashboard

import (
	"github.com/revel/revel"
)

// DashboardController controller
type DashboardController struct {
	DBaseController
}

// Index page
func (c DashboardController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	c.ViewArgs["user"] = c.CurrentUser

	return c.Render()
}

// N make map with size
func N(size int) map[string]interface{} {
	return make(map[string]interface{}, size)
}
