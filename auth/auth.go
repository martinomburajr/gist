package auth

import (
	"encoding/json"
	"fmt"
	"github.com/martinomburajr/gist/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	//RedirectURI is the URI the Github OAuth flow redirects to
	RedirectURI = fmt.Sprintf("http://localhost:%d/auth/github/callback", config.PORT)

	//BaseURL is the base URL to perform a login, this URL does not point to anything by itself,
	// it needs to be composed with other information. See AuthURL for the full Login URL
	BaseURL = "https://github.com/login/oauth/authorize"

	//ClientID represents the client id - this should never be placed in code but rather injected via a variable
	ClientID = ""

	//ClientSecret represents the application client secret - THIS SHOULD NEVER BE PLACED IN CODE BUT INJECTED VIA
	// ENVIRONMENT VARIABLE
	ClientSecret = ""

	//AuthURL is the final URL used to perform a login
	AuthURL = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", BaseURL, ClientID, RedirectURI)

	//Session is a singleton variable that holds all authentication and config based information for a session to
	// succeed.
	Session = SessionObj{}
)

//CreateOAuth2AuthorizationRequest initiates the OAuth2 process. It contacts the GitHub authorization server and requests a code (authorization key) that will later be exchanged for an AccessToken
func CreateOAuth2AuthorizationRequest() error {
	request, err := http.NewRequest(http.MethodGet, AuthURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Print(string(data))

	return nil
}

//RedirectHandler is hit after the CreateAuth2AuthorizationRequest begins the OAuth2 Transaction.
//The response should contain
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("I AM HERE REDIRECTED")
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	code := r.FormValue("code")

	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", ClientID, ClientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not retrieve http request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	req.Header.Set(http.CanonicalHeaderKey("accept"), "application/json")
	req.Header.Set("X-OAuth-Scopes", "gists")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	defer res.Body.Close()
	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	Session.AccessToken = t.AccessToken

	w.WriteHeader(http.StatusFound)
	w.Write([]byte("OK"))
}

//SessionObj is a type that contains session based information for authentication based actions to work
type SessionObj struct {
	AccessToken string `json:"access_token"`
	Client *http.Client
}

//OAuthAccessResponse embodies a response from the GitHub OAuth server with the AccessToken if authorized.
type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
