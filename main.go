package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    httpHandler "gogameserver/httpreqhandler"
)


var (
    AppLogFileName  = os.Getenv("APP_LOG_FILE_PATH")
    GET_TOP_PLAYERS_URL = os.Getenv("GET_TOP_PLAYERS_URL")
    UPDATE_PLAYER_DATA_URL = os.Getenv("UPDATE_PLAYER_DATA_URL")
    GET_PLAYER_RANK_URL = os.Getenv("GET_PLAYER_RANK_URL")
    UPDATE_PLAYER_SCORE_URL = os.Getenv("UPDATE_PLAYER_SCORE_URL")  
    PING_URL = os.Getenv("PING_URL")  
)

func cleaner(){
    for range time.Tick(time.Day *1){
        // delete 
    }
}

func main() {
    go cleaner()
    InitLogging(AppLogFileName)
    
    defineHandlers()

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePing(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Gogameserver")
}

func defineHandlers(){
    http.HandleFun(PING_URL,                 handlePing)
    http.HandleFun(GET_TOP_PLAYERS_URL,     httpHandler.HandleGetTopScorers)
    http.HandleFun(UPDATE_PLAYER_DATA_URL,     httpHandler.HandleUpdatePlayerData)
    http.HandleFun(GET_PLAYER_RANK_URL,     httpHandler.HandleGetPlayerRank)
    http.HandleFun(UPDATE_PLAYER_SCORE_URL,     httpHandler.HandleUpdatePlayerScore )
}

func InitLogging(logFileName string) {
    newLog, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println(err)
    }

    log.SetFlags(0)
    log.SetOutput(newLog)
}
