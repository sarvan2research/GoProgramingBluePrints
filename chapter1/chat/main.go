package main

import (
	"GOProgrammingBluePrints/chapter1/trace"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	data := map[string]interface{}{
		"Host": request.Host,
	}
	if authCookie, err := request.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(writer, data)
}

var host = flag.String("host", ":8080", "The addr of the  application.")

func main() {

	flag.Parse() // parse the flags
	gomniauth.SetSecurityKey("OWN TEST SIGNATURE SARVAN")
	gomniauth.WithProviders(
		facebook.New("930805643740-4ti2p9nplh3flth438fl2711lu4egci1.apps.googleusercontent.com", "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",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:8080/auth/callback/github"),
		google.New("930805643740-4ti2p9nplh3flth438fl2711lu4egci1.apps.googleusercontent.com", "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",
			"http://localhost:8080/auth/callback/google"),
	)
	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	//http.Handle("/", &templateHandler{fileName: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/upload", &templateHandler{fileName: "upload.html"})
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/logout", func(writer http.ResponseWriter, request *http.Request) {
		http.SetCookie(writer, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		writer.Header().Set("Location", "/chat")
		writer.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
