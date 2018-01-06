package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/montera82/go-chat-app/pkg/trace"
	"github.com/stretchr/objx"
)

func main() {

	addrr := flag.String("addr", ":8083", "Specify the host port of the server")
	flag.Parse()

	room := newRoom(useGravatarAvatar)
	room.trace = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateLoader{fileName: "chat.html"}))
	http.Handle("/login", &templateLoader{fileName: "login.html"})
	http.HandleFunc("/auth/", handleAuthenticationProcess)
	http.Handle("/room", room)
	http.HandleFunc("/logout", handleLogout)
	http.Handle("/upload", &templateLoader{fileName: "upload_photo.html"})

	go room.listenForJoinsDepartureAndMessages() //activate the room to listenForJoinsDepartureAndMessages for activities infinitivly

	fmt.Println("Server will runing on port ", *addrr)
	if err := http.ListenAndServe(*addrr, nil); err != nil {
		log.Fatal("Listen and serve", err)
	}

}

type templateLoader struct {
	load             sync.Once
	fileName         string
	compiledTemplate *template.Template
}

func (t *templateLoader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.load.Do(func() {
		t.compiledTemplate = template.Must(template.ParseFiles(filepath.Join("templates/" + t.fileName)))
	})

	dataForChatTemplate := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		dataForChatTemplate["userData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.compiledTemplate.Execute(w, dataForChatTemplate)
}
