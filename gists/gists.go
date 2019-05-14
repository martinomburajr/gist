package gists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/martinomburajr/gist/auth"
	"io/ioutil"
	"net/http"
	"time"
)


// GistFile represents an application facing Gist that a user can create.
// Typically populated through the use of flags. It contains the barebones for what a gist on GitHub may be.
// A GistFile implements a cruder interface and can perform all basic operations.
type GistFile struct {
	Description string         `json:"description"`
	Public      bool           `json:"public"`
	Files       []GistFileBody `json:"file"`
}

//Removes the remote Gist
func (g *GistFile) Delete(id string) (*http.Response, error) {
	urll := fmt.Sprintf("/gists/%s", id)
	req, err := http.NewRequest(http.MethodDelete, urll, nil)
	if err != nil {
		return nil, err
	}

	resp, err := auth.Session.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GistFile) Update(interface{}) (*http.Response, error) {
	panic("implement me")
}

//Given a GistFile in its basic form, create a gist on Github that takes the contents of the Files, description and whether or not it is public.
func (g *GistFile) Create() (*http.Response, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}

	urll := EndpointBase+"/gists"

	req, err := http.NewRequest(http.MethodPost, urll,  bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	resp, err := auth.Session.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//https://developer.github.com/v3/gists/#get-a-single-gist
func  (g *GistFile) Retrieve(id string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, "/gists/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := auth.Session.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	var gf httpGistResponse
	err = json.Unmarshal(data, gf)
	if err != nil {
		return resp, err
	}

	g.Description = gf.Description
	g.Files[0].Content = gf.Files.Gists[0].Content
	g.Public = gf.Public

	return resp, nil
}


type GistFileBody struct {
	FileName string
	Content string `json:"content"`
}

type GistOwner struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
}

type User struct{
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

//returned when a call to retrieve with a gist id is provided.
type httpGistResponse struct {
	URL        string `json:"url"`
	ForksURL   string `json:"forks_url"`
	CommitsURL string `json:"commits_url"`
	ID         string `json:"id"`
	NodeID     string `json:"node_id"`
	GitPullURL string `json:"git_pull_url"`
	GitPushURL string `json:"git_push_url"`
	HTMLURL    string `json:"html_url"`
	Files      struct {
		Gists []httpGistFileResponse
	} `json:"files"`
	Public      bool        `json:"public"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Description string      `json:"description"`
	Comments    int         `json:"comments"`
	User        interface{} `json:"user"`
	CommentsURL string      `json:"comments_url"`
	Owner       GistOwner `json:"owner"`
	Truncated bool `json:"truncated"`
	Forks     []struct {
		User User  `json:"user"`
		URL       string    `json:"url"`
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"forks"`
	History []struct {
		URL     string `json:"url"`
		Version string `json:"version"`
		User    struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		ChangeStatus struct {
			Deletions int `json:"deletions"`
			Additions int `json:"additions"`
			Total     int `json:"total"`
		} `json:"change_status"`
		CommittedAt time.Time `json:"committed_at"`
	} `json:"history"`
}

type httpGistFileResponse struct {
	Filename  string `json:"filename"`
	Type      string `json:"type"`
	Language  string `json:"language"`
	RawURL    string `json:"raw_url"`
	Size      int    `json:"size"`
	Truncated bool   `json:"truncated"`
	Content   string `json:"content"`
}