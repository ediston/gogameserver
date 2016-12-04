package main

import (
    "fmt"
    "log"
    "net/http"

    rcl "gogameserver/redisclient" 
)

func echoString(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello")
}

func main() {
    rc := rcl.New()
    keyStr := "00NeverAddThiskey"
    valStr := "00NeverAddThisVal"
    rc.SaveKeyValForever(keyStr, valStr)
    tempVal, _ := rc.GetVal(keyStr)
    if valStr != tempVal{
        fmt.Sprintf("Key should exist and be equal to %s!", valStr)
    } else {
        fmt.Sprintf("valStr = =%s", valStr)
    }
    rc.DelKey(keyStr)

    http.HandleFunc("/", echoString)
    log.Fatal(http.ListenAndServe(":8081", nil))

}