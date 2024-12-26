package main

import apiserver "github.com/summerliuguang/learngo/apiserver"

// func init() {
//     file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//     if err!= nil {
//         log.Fatal(err)
//     }
//     log.SetOutput(file)
// }

func main() {
	var server = apiserver.NewAPIServer(":8080")
	server.Run()
}
