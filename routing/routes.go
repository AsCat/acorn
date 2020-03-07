package routing

import (
	"net/http"

	"github.com/AsCat/acorn/handlers"
)

// Route describes a single route
type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	Authenticated bool
}

// Routes holds an array of Route. A note on swagger documentation. The path variables and query parameters
// are defined in ../doc.go.  YOu need to manually associate params and routes.
type Routes struct {
	Routes []Route
}

// NewRoutes creates and returns all the API routes
func NewRoutes() (r *Routes) {
	r = new(Routes)

	r.Routes = []Route{
		// swagger:route GET / Root
		// ---
		// Endpoint to get the health of Kiali
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		// responses:
		//      200: statusInfo
		{
			"Healthz",
			"GET",
			"/healthz",
			handlers.Healthz,
			false,
		},
		// swagger:route GET / Root
		// ---
		// Endpoint to get the status of Kiali
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		// responses:
		//      500: internalError
		//      200: statusInfo
		{
			"Root",
			"GET",
			"/api",
			handlers.Root,
			false,
		},

		{
			"NamespaceList",
			"GET",
			"/api/namespaces",
			handlers.NamespaceList,
			true,
		},
	}

	return
}
