package restroute_test

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ceralena/go-restroute"
)

type uc interface {
	GetUsers() []interface{}
	GetUser(user_id string) interface{}
	CreateUser(json_body io.Reader) interface{}
	UpdateUser(user_id string, json_body io.Reader) interface{}
}

var userController uc

func ExampleMap() {
	// define our route handlers
	getUsers := func(req restroute.Request) {
		users := userController.GetUsers()
		enc := json.NewEncoder(req.W)
		enc.Encode(users)
	}

	getUser := func(req restroute.Request) {
		user_id := req.Params["user_id"]
		users := userController.GetUser(user_id)
		enc := json.NewEncoder(req.W)
		enc.Encode(users)
	}

	createUser := func(req restroute.Request) {
		users := userController.CreateUser(req.R.Body)
		enc := json.NewEncoder(req.W)
		enc.Encode(users)
	}

	updateUser := func(req restroute.Request) {
		user_id := req.Params["user_id"]
		users := userController.UpdateUser(user_id, req.R.Body)
		enc := json.NewEncoder(req.W)
		enc.Encode(users)
	}

	// map RESTful route methods to the handlers
	routes := restroute.Map{
		`/users`: {
			"GET": getUsers,
			"PUT": createUser,
		},
		`/users/(?P<user_id>\w+)`: {
			"GET":  getUser,
			"POST": updateUser,
		},
	}

	// listen for http over port 8080 with our route map as the handler
	http.ListenAndServe(":8080", routes.MustCompile())
}

// package restroute does not have any concept of a middleware as such - but
// you take advantage of function objects and closures to build middleware
// very easily.
//
// This approach makes it fairly easy to compose a chain of
// middlewares to handle strongly typed dependency injection.
//
// This example implements the same users router as above, but with a
// middleware so that routers don't have to handle JSON encoding.
//
func ExampleMap_middleware() {
	// handleJSON is our middleware - it takes a function that receives a
	// request and returns some JSON, and returns an ordinary handler which
	// can be used in restroute.Map.
	handleJSON := func(handler func(req restroute.Request) interface{}) func(req restroute.Request) {
		fn := func(req restroute.Request) {
			data := handler(req)
			enc := json.NewEncoder(req.W)
			enc.Encode(data)
		}
		return fn
	}

	// define our route handlers
	getUsers := func(req restroute.Request) interface{} {
		return userController.GetUsers()
	}

	getUser := func(req restroute.Request) interface{} {
		user_id := req.Params["user_id"]
		return userController.GetUser(user_id)
	}

	createUser := func(req restroute.Request) interface{} {
		return userController.CreateUser(req.R.Body)
	}

	updateUser := func(req restroute.Request) interface{} {
		user_id := req.Params["user_id"]
		return userController.UpdateUser(user_id, req.R.Body)
	}

	// map RESTful route methods to the handlers
	routes := restroute.Map{
		`/users`: {
			"GET": handleJSON(getUsers),
			"PUT": handleJSON(createUser),
		},
		`/users/(?P<user_id>\w+)`: {
			"GET":  handleJSON(getUser),
			"POST": handleJSON(updateUser),
		},
	}

	// listen for http over port 8080 with our route map as the handler
	http.ListenAndServe(":8080", routes.MustCompile())
}
