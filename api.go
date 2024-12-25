package main

import (
	"log"
	"net/http"
)

type APIServer struct {
	address string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		address: addr,
	}
}

func (s *APIServer) Run() error {
	route := http.NewServeMux()
	route.HandleFunc("/getuser/{user}", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		w.Write([]byte("Hello, " + user))
	})
	server := http.Server{
		Addr:    s.address,
		Handler: route,
	}
	log.Println("Server started at", s.address)
	err := server.ListenAndServe()
	return err
}
