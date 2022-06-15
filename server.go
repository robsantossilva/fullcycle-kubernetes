package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/configmap", ConfigMap)
	http.HandleFunc("/", Hello)
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

func Hello(w http.ResponseWriter, r *http.Request) {

	name := os.Getenv("NAME")
	age := os.Getenv("AGE")

	fmt.Fprintf(w, "Hello %s, you are %s years old", name, age)

	//w.Write([]byte("<h1>Hello, world!</h1>"))
}

func ConfigMap(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadFile("myfamily/family.txt")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	fmt.Fprintf(w, "My Family: %s", string(data))
}
