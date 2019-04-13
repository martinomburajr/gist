package gists

import "net/http"

const (
	EndpointBase = "https://api.github.com/"
	EndpointGistCreate = "/gists"
	EndpointGistCreateMethod = http.MethodPost
)