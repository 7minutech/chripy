package main

import "net/http"

var serverMux = http.ServeMux{}

var localServer = http.Server{
	Handler: &serverMux,
	Addr:    ":8080",
}
