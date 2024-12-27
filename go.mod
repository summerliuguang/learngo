module github.com/summerliuguang/learngo

go 1.23.4

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.31.0 // indirect
)

replace (
	github.com/summerliuguang/learngo/apiserver => ./apiserver
	github.com/summerliuguang/learngo/pqcontrol => ./pqcontrol
)
