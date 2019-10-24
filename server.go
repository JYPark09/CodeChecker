package main

import (
	"log"
	"net/http"
)

func startServer(port string) *http.Server {
	srv := &http.Server{Addr: port}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/run", runHandler)
	http.HandleFunc("/result", resultHandler)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("[http] Listen failed ", err)
		}
	}()

	return srv
}
