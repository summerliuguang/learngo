package apiserver

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	pqcontrol "github.com/summerliuguang/learngo/pqcontrol"
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
	router.HandleFunc("/getuser/{userid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userid := vars["userid"]
		name, result := pqcontrol.GetUserById(userid)
		if result != pqcontrol.Success {
			http.Error(w, "Get user failed", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Hello, " + name + "\n"))
	}).Methods("GET")

	router.HandleFunc("/getuserlist", func(w http.ResponseWriter, r *http.Request) {
		users, result := pqcontrol.GetUsers()
		if result != pqcontrol.Success {
			http.Error(w, "Get users failed", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Users: " + strings.Join(users, ", ") + "\n"))

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
