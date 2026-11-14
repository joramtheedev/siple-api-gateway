package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "9001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "backend :%s got %s %s\n", port, r.Method, r.URL.Path)
	})
	log.Println("mock backend listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
