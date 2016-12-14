package httpreqhandler

import (
    "bytes"
    "net/http"
    "log"
    "time"
    "encoding/json"

    util "gogameserver/util"
    gdm "gogameserver/gamedatamanager"
    dt  "gogameserver/datatypes" 
)

const WEEK string  = "w" 
const DAY string  = "d" 
const EVER string  = "e"
const GAMENAMESHORT string = "gn"
const TYPESHORT string = "ty"
const PLAYERIDSHORT string = "p"
const TOPAMOUNTSHORT string = "ta"
const SCORESHORT string = "s"
const PLAYERRANK string = "pr"
const RESP string  = "resp"

func handleEmptyInputs(respWriter http.ResponseWriter, reqFromClient *http.Request, gameName string, playerId string){
    respToClient := map[string]interface{}{
                GAMENAMESHORT : gameName,
                PLAYERIDSHORT : playerId,
            }
    respData := sendResponse(respWriter, reqFromClient, respToClient)
    go log.Printf("Error: GameName: (%s) or playerId: (%s) is empty. respData: %s", gameName, playerId, respData)
}

func HandleGetTopScorers(respWriter http.ResponseWriter, reqFromClient *http.Request) {
    var respToClient map[string]interface{}

    start := time.Now()

    topAmount := util.URLParamInt(reqFromClient, TOPAMOUNTSHORT, 10)  // top amount
    topType   := util.URLParamStr(reqFromClient, TYPESHORT, EVER)  // type "week", "day", "ever"
    gameName  := util.URLParamStr(reqFromClient, GAMENAMESHORT, "") // gamename:
    playerId  := util.URLParamStr(reqFromClient, PLAYERIDSHORT, "")  //  player id
    
    if gameName == "" || playerId == "" {
        handleEmptyInputs(respWriter, reqFromClient, gameName, playerId)
        return
    }

    gm := gdm.New() // game manager

    respToClient = map[string]interface{}{
        PLAYERIDSHORT : playerId,
        GAMENAMESHORT : gameName,
        WEEK : "",
        DAY : "",
        EVER : "",
    }

    switch topType {
        case WEEK: 
            respToClient[WEEK] = gm.GetTopPlayersThisWeek(gameName , topAmount)
        
        case DAY: 
            respToClient[DAY] = gm.GetTopPlayersOnDay(gameName , topAmount, 0)
        
        case EVER: 
            respToClient[EVER] = gm.GetTopPlayers(gameName, topAmount)
    }

    respData := sendResponse(respWriter, reqFromClient, respToClient)
    elapsed := time.Since(start) / time.Millisecond
    go logResponse(respData, elapsed)
}

func HandleGetPlayerRank(respWriter http.ResponseWriter, reqFromClient *http.Request) {
    var respToClient map[string]interface{}

    start := time.Now()

    gameName  := util.URLParamStr(reqFromClient, GAMENAMESHORT, "") // gamename:
    playerId := util.URLParamStr(reqFromClient, PLAYERIDSHORT, "")  //  
    
    if gameName == "" || playerId == "" {
        handleEmptyInputs(respWriter, reqFromClient, gameName, playerId)
        return
    }

    gm := gdm.New() // game manager
    
    respToClient = map[string]interface{}{
        PLAYERRANK    : gm.GetPlayerRank(gameName , playerId ),
        PLAYERIDSHORT : playerId,
        GAMENAMESHORT : gameName,
    }

    respData := sendResponse(respWriter, reqFromClient, respToClient)
    elapsed := time.Since(start) / time.Millisecond
    go logResponse(respData, elapsed)
}

func HandleUpdatePlayerData(respWriter http.ResponseWriter, reqFromClient *http.Request) {
    var respToClient map[string]interface{}
    var playerDataJson dt.PlayerData
    start := time.Now()

    gameName  := util.URLParamStr(reqFromClient, GAMENAMESHORT, "") // gamename:
    playerId := util.URLParamStr(reqFromClient, PLAYERIDSHORT, "")  //  player id

    if reqFromClient.Body == nil || gameName == "" || playerId == "" {
        handleEmptyInputs(respWriter, reqFromClient, gameName, playerId)
        return
    }

    err := json.NewDecoder(reqFromClient.Body).Decode(&playerDataJson)
    if err != nil {
        handleEmptyInputs(respWriter, reqFromClient, gameName, playerId + ": error")
        return
    }

    gm := gdm.New() // game manager

    respToClient = map[string]interface{}{
        RESP : gm.StorePlayerData(gameName, playerDataJson),
    }

    respData := sendResponse(respWriter, reqFromClient, respToClient)
    elapsed := time.Since(start) / time.Millisecond
    go logResponse(respData, elapsed)
}

func HandleUpdatePlayerScore(respWriter http.ResponseWriter, reqFromClient *http.Request) {
    var respToClient map[string]interface{}

    start := time.Now()

    gameName  := util.URLParamStr(reqFromClient, GAMENAMESHORT, "") // gamename:
    playerId := util.URLParamStr(reqFromClient, PLAYERIDSHORT, "")  //  player id
    score := util.URLParamFloat(reqFromClient, SCORESHORT, 0)  //  score

    if gameName == "" || playerId == "" {
        handleEmptyInputs(respWriter, reqFromClient, gameName, playerId)
        return
    }

    gm := gdm.New() // game manager

    respToClient = map[string]interface{}{
        RESP : gm.StorePlayerScore(gameName, score, playerId),
    }

    respData := sendResponse(respWriter, reqFromClient, respToClient)
    elapsed  := time.Since(start) / time.Millisecond
    go logResponse(respData, elapsed)
}

func getJsonResp(respToClient map[string]interface{}) []byte {
    jsonResp, err := json.MarshalIndent(respToClient, "", " ")
    if err != nil {
        log.Print(err)
    }
    return jsonResp
}

func sendResponse(respWriter http.ResponseWriter, httpReq *http.Request, respToClient map[string]interface{}) []byte {
    jsonResp, err := json.MarshalIndent(respToClient, "", " ")
    if err != nil {
        log.Print(err)
    }
    callback := util.URLParamStr(httpReq, "callback", "")
    respWriter.Header().Set("Expires", "-1")
    respWriter.Header().Set("Cache-Control", "private, max-age=0")
    if callback == "" {
        respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
    } else {
        respWriter.Header().Set("Content-Type", "application/javascript")
        cb := []byte(callback)
        jsonResp = bytes.Join([][]byte{cb, []byte("("), jsonResp, []byte(")")}, []byte{})
    }
    respWriter.Write(jsonResp)
    return jsonResp
}

func logResponse(data []byte, elapsed time.Duration) { 
    logMsg := "INFO: Response time: [%v] \"GET /api/v1/results HTTP/1.0\" %v %v %v"
    millisec := (float64(elapsed) / float64(time.Millisecond)) * 1000000
    log.Printf("Completed request in %v ms", millisec)
    log.Printf(logMsg, time.Now(), "200", len(data), millisec*1000)
}
