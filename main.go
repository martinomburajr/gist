package main

import (
	"fmt"
	mux2 "github.com/gorilla/mux"
	"github.com/martinomburajr/gist/auth"
	"github.com/martinomburajr/gist/config"
	"html/template"
	"log"
	"net/http"
)

func main() {
	//flags
	//file := flag.String("file", "", "file to be sent to gist")
	//fileOverride := flag.String("setfile", "", "overrides the current file name for a new one. The filename provided by this flag will be reflected on the gists.github.com website")
	//description := flag.String("description", "", "sets the description of the gist")
	//isPublic := flag.Bool("pub", true, "set as public gist. This is set to true by default")

	//if file != "" {
	//	//check file exists
	//}

	//@todo change mux2 alias to original mux alias
	mux := mux2.NewRouter()

	authTemplate := template.Must(template.ParseFiles("public/oauth.html"))
	mux.Methods(http.MethodGet).Path("/").HandlerFunc(LoginHandler(authTemplate))
	mux.Methods(http.MethodGet).Path("/auth/github/callback").HandlerFunc(auth.RedirectHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.PORT), mux))
}

// LoginHandler handles the logging in of a user.
// It will open a simple OAuth Page on a browser that will enable the OAuth flow to begin.
// A successful login returns a valid OAuth AccessToken that is stored in the auth.Session variable.
func LoginHandler(authTemplate *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := authTemplate.Execute(w, struct {
			URLL string
		}{
			URLL: auth.AuthURL,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}


// baseHandler is bare test handler.
// @todo remove baseHandler in main.go if it is expendable
func baseHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello User!")
	err := auth.CreateOAuth2AuthorizationRequest()

	if err != nil {
		return
	}
}