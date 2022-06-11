package main

import "net/http"

func main() {
	http.HandleFunc("/", Hello)
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, world!</h1>"))
}
