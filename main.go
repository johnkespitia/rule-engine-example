package main

import (
    "fmt"
    "log"
    "net/http"
) 

func main() {

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w,"nothing to show !!")
    })

	fmt.Printf("Starting server at port 89-0\n")
	if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatal(err)
    }
}
