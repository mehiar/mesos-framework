package restapi

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"RequestsIndex",
		"GET",
		"/container-requests",
		RequestsIndex,
	},
	Route{
		"RequestCreate",
		"POST",
		"/container-requests",
		RequestCreate,
	},
	Route{
		"RequestShow",
		"GET",
		"/container-requests/{containerId}",
		RequestShow,
	},
}
