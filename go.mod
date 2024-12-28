module github.com/summerliuguang/learngo

go 1.23.4

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/sony/sonyflake v1.2.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
)

replace (
	github.com/summerliuguang/learngo/apiserver => ./apiserver
	github.com/summerliuguang/learngo/pqcontrol => ./pqcontrol
)
