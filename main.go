package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"miniChat/trace"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct{
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCoolie, err := r.Cookie("auth"); err == nil{
		data["UserData"] = objx.MustFromBase64(authCoolie.Value)
	}
	t.templ.Execute(w, data)
}

func main(){
	gomniauth.SetSecurityKey("zuxt")
	gomniauth.WithProviders(
		google.New("354757645848-pb8k2vbh4hhh83gu120bf7shj8i1qp9p.apps.googleusercontent.com", "EVx0Gsq5mTWeEeUdjy8G8k_W", "http://localhost:8080/auth/callback/google"),
		)
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	//r := newRoom()
	r := newRoom(UserGravatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename:"chat.html"}))
	http.Handle("/room", r)
	http.Handle("/login", &templateHandler{filename:"login.html"})
	http.Handle("/upload", &templateHandler{filename:"upload.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/avatars",http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))


	go r.run()

	log.Println("Webサーバーを開始します。ポート：", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil{
		log.Fatal("ListenAndServe:", err)
	}
}