module github.com/summerliuguang/learngo

go 1.23.4

require (
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace (
	github.com/summerliuguang/learngo/apiserver => ./apiserver
	github.com/summerliuguang/learngo/pqcontrol => ./pqcontrol
)