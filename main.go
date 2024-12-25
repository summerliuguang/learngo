package main

func main() {
	var server = NewAPIServer(":8080")
	server.Run()
}
