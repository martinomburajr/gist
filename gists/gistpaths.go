package gists

import "net/http"

const (
	//EndpointBase is the Github API endpoint base URL that can be concatenated with other paths to access more
	// resources
	EndpointBase = "https://api.github.com/"

	//EndpointGistCreate refers to the gist endpoint after being concatenated with EndpointBase
	EndpointGistCreate = "gists"

	//EndpointGistCreateMethod is the appropriate HTTP method for creating a gist
	EndpointGistCreateMethod = http.MethodPost
)