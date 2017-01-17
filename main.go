package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    httpHandler "gogameserver/httpreqhandler"
)


var (
    AppLogFileName  = os.Getenv("APP_LOG_FILE_PATH")
    GET_TOP_PLAYERS_URL = os.Getenv("GET_TOP_PLAYERS_URL")
    UPDATE_PLAYER_DATA_URL = os.Getenv("UPDATE_PLAYER_DATA_URL")
    GET_PLAYER_RANK_URL = os.Getenv("GET_PLAYER_RANK_URL")
    UPDATE_PLAYER_SCORE_URL = os.Getenv("UPDATE_PLAYER_SCORE_URL")  
    UPDATE_PLAYER_DATA_SCORE_AND_GET_RANK_URL = os.Getenv("UPDATE_PLAYER_DATA_SCORE_AND_GET_RANK_URL")
    PING_URL = os.Getenv("PING_URL")  
)

func cleaner(){
    for range time.Tick(60*time.Minute){
        // delete 
    }
}

func main() {
    //go cleaner()
    InitLogging(AppLogFileName)
    
    defineHandlers()

    http.ListenAndServe(":8080", nil)
}

func handlePing(w http.ResponseWriter, req *http.Request) {
    log.Printf("Got Ping Request.\n", )
    fmt.Fprintf(w, "Gogameserver")
}

func defineHandlers(){
    http.HandleFunc(PING_URL,                 handlePing)
    http.HandleFunc(GET_TOP_PLAYERS_URL,     httpHandler.HandleGetTopScorers)
    http.HandleFunc(UPDATE_PLAYER_DATA_URL,     httpHandler.HandleUpdatePlayerData)
    http.HandleFunc(GET_PLAYER_RANK_URL,     httpHandler.HandleGetPlayerRank)
    http.HandleFunc(UPDATE_PLAYER_SCORE_URL,     httpHandler.HandleUpdatePlayerScore )
    http.HandleFunc(UPDATE_PLAYER_DATA_SCORE_AND_GET_RANK_URL,
                    httpHandler.HandleUpdatePlayerDataWithGetPlayerRank )
}

func InitLogging(logFileName string) {
    newLog, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println(err)
    }

    log.SetFlags(0)
    log.SetOutput(newLog)
}
