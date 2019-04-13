package auth

import (
	"encoding/json"
	"fmt"
	"github.com/martinomburajr/gogist/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	RedirectURI = fmt.Sprintf("http://localhost:%d/auth/github/callback", config.PORT)
	BaseURL = "https://github.com/login/oauth/authorize"
	ClientID = ""
	ClientSecret = ""
	AuthURL = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", BaseURL, ClientID, RedirectURI)
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
		fmt.Fprint(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	code := r.FormValue("code")

	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", ClientID, ClientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not retrieve http request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	req.Header.Set(http.CanonicalHeaderKey("accept"), "application/json")
	req.Header.Set("X-OAuth-Scopes", "gists")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	defer res.Body.Close()
	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	Session.AccessToken = t.AccessToken

	w.WriteHeader(http.StatusFound)
	w.Write([]byte("OK"))
}

type SessionObj struct {
	AccessToken string `json:"access_token"`
	Client *http.Client
}
type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
