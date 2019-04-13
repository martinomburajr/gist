package http

import (
	"net/http"
)

//A GistCruder is a generic interface that is implemented by objects the perform Create Retrieve Update and Delete (CRUD) operations.
type GistCruder interface {
	//Performs a create operation using a given id. Typically remote calls on a REST API return HTTP Status Code 201 - Created.
	Create() (*http.Response, error)

	//Performs a get operation using a given id. Typically remote calls on a REST API return HTTP Status Code 200 - OK.
	Retrieve(id string) (*http.Response, error)

	//Performs a delete operation using a given id. Typically remote calls on a REST API return HTTP Status Code 204 - No Content.
	// However depending on the implementation of the delete action, one would need to be aware of the appropriate response
	Delete(id string) (*http.Response, error)

	//Given an interface, the Update function will attempt to Swap out the oldObject with the newObj. Typically remote calls on a REST API return HTTP Status Code 200 - OK.
	Update(newObj interface {}) (*http.Response, error)
}
