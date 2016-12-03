package main

import (
    "fmt"
    "log"
    "net/http"

)

func echoString(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello")
}

func main() {
    http.HandleFunc("/", echoString)
    log.Fatal(http.ListenAndServe(":8081", nil))
}