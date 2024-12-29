package apiserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	pqcontrol "github.com/summerliuguang/learngo/pqcontrol"
)

type APIServer struct {
	address string
}

// Register Routers
func (s *APIServer) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/getuser/{userid}", getUserById).Methods("GET")
	router.HandleFunc("/getuserlists", getUserList).Methods("GET")
	router.HandleFunc("/login", loginAuthentication).Methods("POST")
	router.HandleFunc("/register", registerUser).Methods("POST")
}

// Register Middleware
func (s *APIServer) registerMiddlewareV1(router *mux.Router) {
	router.Use(requestLoggerMiddleware)
	router.Use(requestAuthMiddleware)
}

func (s *APIServer) RegisterMiddlewareCommon(router *mux.Router) {
	router.Use(requestLoggerMiddleware)
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		address: addr,
	}
}
func (s *APIServer) createSubRouter(router *mux.Router, apipath string) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix(apipath).Handler(http.StripPrefix(apipath, router))
	return r
}

func (s *APIServer) configureRouter() *mux.Router {
	// main router
	router := mux.NewRouter()
	s.RegisterRouter(router)

	v1 := s.createSubRouter(router, "/api/v1")
	s.registerMiddlewareV1(v1)
	common := s.createSubRouter(router, "/common")
	s.RegisterMiddlewareCommon(common)

	r := mux.NewRouter()
	r.PathPrefix("/api/v1").Handler(v1)
	r.PathPrefix("/common").Handler(common)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	r.Use(func(next http.Handler) http.Handler {
		return corsHandler.Handler(next)
	})

	return r
}

func (s *APIServer) Run() error {

	r := s.configureRouter()

	server := http.Server{
		Addr:    s.address,
		Handler: r,
	}
	log.Println("Server started at", s.address)
	err := server.ListenAndServe()
	return err
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("From:%s, Method:%s, URL:%s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func requestAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := ValidJWT(token)
		if err != nil {
			http.Error(w, "Invalid JWt", http.StatusUnauthorized)
			return
		}
		type contextKey string
		const authUsernameKey contextKey = "auth_username"
		ctx := context.WithValue(r.Context(), authUsernameKey, claims.Username)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["userid"]
	username := vars["auth_username"]
	name, result := pqcontrol.GetUserById(userid)
	if result != pqcontrol.Success {
		http.Error(w, "Get user failed", http.StatusInternalServerError)
		return
	}
	log.Println("auth:", username)
	w.Write([]byte("Hello, " + name + "\n"))
}

func getUserList(w http.ResponseWriter, r *http.Request) {
	users, result := pqcontrol.GetUsers()
	if result != pqcontrol.Success {
		http.Error(w, "Get users failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Users: " + strings.Join(users, ", ") + "\n"))
}

func loginAuthentication(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	username, result := pqcontrol.AuthAccount(loginRequest.Username, loginRequest.Password)
	if result != pqcontrol.Success {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}
	signedToken, err := GenerateJWT(username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+signedToken)
	w.Write([]byte(username + " login successful\n"))
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// check turnstile token in request header
	// if r.Method == http.MethodPost {
	// 	fmt.Println(r.Header)
	// }
	var turnstileVerify TurnstileVerify
	turnstileVerify.Response = loginRequest.TurnstileToken
	turnstileStatus := VaildTurnstile(turnstileVerify)
	if !turnstileStatus {
		http.Error(w, "Invalid Turnstile Token", http.StatusUnauthorized)
		return
	}

	userid, result := pqcontrol.CreateAccount(loginRequest.Username, loginRequest.Password)
	response := Response{
		Code:    0,
		Message: "",
	}
	w.Header().Set("Content-Type", "application/json")
	if result != pqcontrol.Success {
		if result == pqcontrol.UserAlreadyExists {
			response.Code = pqcontrol.UserAlreadyExists
			response.Message = "Error: User already exists"
			responseData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Error marshalling response", http.StatusInternalServerError)
				return
			}
			w.Write(responseData)
		} else {
			http.Error(w, "User creation failed", http.StatusInternalServerError)
		}
		return
	}
	response.Code = pqcontrol.Success
	response.Message = "User: " + strconv.FormatInt(userid, 10) + " created successfully"
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Write(responseData)
}
