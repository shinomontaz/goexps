package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	port int
	host string
)

func init() {
	pp := flag.Int("port", 8080, "port")
	flag.Parse()

	port = *pp

	var err error

	host, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", simple)

	fmt.Printf("start listening on port :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func simple(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hit %s\n", host)
}
