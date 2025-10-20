package main

import "net/http"

const port = "8080"

var serverMux = http.NewServeMux()

var localServer = http.Server{
	Addr:    ":" + port,
	Handler: serverMux,
}
