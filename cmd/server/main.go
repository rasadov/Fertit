package main

import "net/http"

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", fileServer)
}
