package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandlerStruct struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (t *templateHandlerStruct) ServeHttp(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	t.templ.Execute(w, nil)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`<html> 
        <head> 
          <title>Chat</title> 
        </head> 
        <body> 
          Let's chat! 
        </body> 
      </html>`))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
