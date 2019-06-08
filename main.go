package main

import (
	"log"
	"net/http"
	"sync"
	"text/template"
)

type templateHandler struct{
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	t.once.Do()
}

func main(){

	if err := http.ListenAndServe(":8080", nil); err != nil{
		log.Fatal("ListenAndServe:", err)
	}
}