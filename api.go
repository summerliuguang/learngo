package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	address string
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(m ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next.ServeHTTP
	}
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		address: addr,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/getuser/{user}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		log.Println("User:", vars)
		w.Write([]byte("Hello, " + user))
	}).Methods("GET")

	v1 := mux.NewRouter()
	v1.PathPrefix("/api/v1").Handler(http.StripPrefix("/api/v1", router))

	middlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		RequestAuthMiddleware,
	)

	server := http.Server{
		Addr:    s.address,
		Handler: middlewareChain(v1),
	}
	log.Println("Server started at", s.address)
	err := server.ListenAndServe()
	return err
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("From:%s, Method:%s, URL:%s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	}
}

func RequestAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
