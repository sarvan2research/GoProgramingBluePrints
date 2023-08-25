package main

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/gomniauth"
	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"strings"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}
type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if cookie, err := request.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated by checking auth cookie not found
		writer.Header().Set("Location", "/login")
		writer.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// different error getting cookie
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	} else {
		// success call
		a.next.ServeHTTP(writer, request)
	}
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
		chatUser := &chatUser{
			User: user,
		}
		m := md5.New()
		io.WriteString(m, strings.ToLower(strings.TrimSpace(user.Email())))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("Error when trying to GetAvatarURL", "-", err)
		}
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
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
