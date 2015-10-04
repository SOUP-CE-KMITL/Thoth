// TODO : need to change to api
package main

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
		"GetNodes",
		"GET",
		"/nodes",
		GetNodes,
	},
	Route{
		"Cors",
		"OPTIONS",
		"/",
		OptionCors,
	},
	Route{
		"GetNode",
		"GET",
		"/node/{nodeName}",
		GetNode,
	},
	Route{
		"GetPods",
		"GET",
		"/pods",
		GetPods,
	},
	Route{
		"GetPod",
		"GET",
		"/pod/{podName}",
		GetPod,
	},
	Route{
		"CreatePod",
		"POST",
		"/pod/create",
		CreatePod,
	},
}
