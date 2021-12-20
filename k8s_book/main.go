package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	port string
	host string
)

func init() {
	//	port = os.Getenv("PORT")
	port = "8080"
	var err error

	host, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", simple)

	fmt.Printf("start listening on port :%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func simple(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hit %s\n", host)
}
