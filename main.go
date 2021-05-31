package main

import (
	"net/http"

	"github.com/erda-project/erda-bot/conf"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda/pkg/httpserver"
)

func main() {
	// load conf
	conf.Load()
	// load gh
	gh.Init()

	// init http server
	server := httpserver.New(":4567")
	routers := []httpserver.Endpoint{
		{Path: "/webhooks", Method: http.MethodPost, Handler: Webhooks},
	}
	server.RegisterEndpoint(routers)
	panic(server.ListenAndServe())
}
