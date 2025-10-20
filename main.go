package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	var localServer = http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(localServer.ListenAndServe())

}
