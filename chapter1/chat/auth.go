package main

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	_, err := request.Cookie("auth")
	if err == http.ErrNoCookie {
		writer.Header().Set("Location", "/login")
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	a.next.ServeHTTP(writer, request)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
// key=API_KEY AIzaSyB0xeSzcp_lCAkSSM1JEk7w1K9opmZ8qV4
func loginHandler(writer http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	log.Println("Length of url segments:", segs)
	if len(segs) <= 2 {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "Url Malformed ")
		return
	}
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		log.Println("TODO handle for login for", provider)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error when trying to GetBeginAuthURLfor %s:%s", provider, err), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Location", loginUrl)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error when trying to get provider %s: %s",
				provider, err), http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error when trying to get user from %s: %s",
				provider, err), http.StatusInternalServerError)
			return
		}
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(writer, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		writer.Header().Set("Location", "/chat")
		writer.WriteHeader(http.StatusTemporaryRedirect)
	default:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "Auth action %s not supported", action)
	}

}
