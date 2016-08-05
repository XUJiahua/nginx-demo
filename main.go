package main

import (
	"log"
	"net/http"
	"flag"
)

type demo struct {

}

func (*demo)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(port))
}

var port = "10000"

func main() {
	p := flag.Bool("p", false, "use 20000")
	flag.Parse()
	if *p == true {
		port = "20000"
	}
	http.Handle("/", &demo{})
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
