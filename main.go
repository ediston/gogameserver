package main

import (
    "fmt"
    "log"
    "net/http"
)


var (
    AppLogFileName      = "/app/log/app.log"
)

func cleaner(){
    for range time.Tick(time.Day *1){
        // delete 
    }
}

func echoString(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello")
}

func main() {
    go cleaner()
    InitLogging(AppLogFileName)
  
    http.HandleFunc("/", echoString)
    log.Fatal(http.ListenAndServe(":8081", nil))
}

func InitLogging(logFileName string) {
    newLog, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println(err)
    }
    
    log.SetFlags(0)
    log.SetOutput(newLog)
}
