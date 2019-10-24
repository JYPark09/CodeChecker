package main

import (
	"html/template"
	"net/http"
)

func mainHandler(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("template/index.html")

	t.Execute(writer, &langMgr)
}
