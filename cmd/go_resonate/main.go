package main

import (
	"github.com/zballs/go_resonate/api"
	. "github.com/zballs/go_resonate/util"
	"net/http"
)

func main() {

	CreatePages(
		"artist",
		"listener",
		"login",
	)

	RegisterTemplates(
		"artist.html",
		"listener.html",
		"login.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/artist", TemplateHandler("artist.html"))
	mux.HandleFunc("/listener", TemplateHandler("listener.html"))
	mux.HandleFunc("/login", TemplateHandler("login.html"))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
