package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	s := 0
	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		s++
		fmt.Printf("[origin server %d]\n", s)
		_, _ = fmt.Fprint(rw, "response from "+os.Args[1])
	})

	log.Fatal(http.ListenAndServe(":"+os.Args[1], originServerHandler))
}
