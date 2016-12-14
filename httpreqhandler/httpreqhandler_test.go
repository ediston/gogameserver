package httpreqhandler

import (
    dt  "gogameserver/datatypes"
    gdm "gogameserver/gamedatamanager"

    "testing"
    "strconv"
    "bytes"
    "net/http"
    "net/http/httptest"
)

const GAMENAME string  = "00dummygame" 
const PLAYERID string  = "00playerid"
var PLAYERIDS  = [] string{"00NeverAddThispid0", "00NeverAddThispid1",  "00NeverAddThispid3", "00NeverAddThispid4", "00NeverAddThispid5"}



func TestHandleUpdatePlayerData(t *testing.T) {
    var url         bytes.Buffer

    gameName    := GAMENAME
    playerId    := PLAYERID
    playerDataStr  := dt.Str(dt.NewWithId(PLAYERID))

    url.WriteString("http://gogameserver.com/results?")
    url.WriteString(GAMENAMESHORT + "="+gameName + "&")
    url.WriteString(PLAYERIDSHORT + "="+playerId + "&") 

    urlStr   := url.String()

    req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(playerDataStr)) )
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("X-Custom-Header", "value")
    req.Header.Set("Content-Type", "application/json")

    // Let game manager delete an existing data with the gamename and playerid
    gm := gdm.New()
    gm.DelPlayerData(GAMENAME, PLAYERID)

    playerDataStrFromDB, found := gm.GetPlayerData(GAMENAME, PLAYERID)
    if found {
        t.Errorf("Error: GetPlayerData: Data shall not have been found\n")
    }  

    // continue
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleUpdatePlayerData)

    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // check using game data if the newly entered data exist
    playerDataStrFromDB, found = gm.GetPlayerData(GAMENAME, PLAYERID)
    if !found {
        t.Errorf("Error: TestHandleUpdatePlayerData: Data shall have been found, playerDataStrFromDB:%s\n", playerDataStrFromDB)
    } else {
        if playerDataStrFromDB != playerDataStr {
            t.Errorf("Error: TestHandleUpdatePlayerData: playerDataStr from redis: %s\n\tplayerDataStr shall have been %s\n", playerDataStrFromDB, playerDataStr)
        }
    }
    // delete finally
    gm.DelPlayerData(GAMENAME, PLAYERID)
}

func TestHandleUpdatePlayerScore(t *testing.T) {
    var url         bytes.Buffer

    gameName    := GAMENAME
    playerId    := PLAYERID
    score       := 20
    scoreF      := 20.0

    url.WriteString("http://gogameserver.com/results?")
    url.WriteString(GAMENAMESHORT + "="+gameName + "&")
    url.WriteString(PLAYERIDSHORT + "="+playerId + "&") 
    url.WriteString(SCORESHORT + "="+ strconv.Itoa(score ) )

    urlStr   := url.String()

    req, _ := http.NewRequest("GET", urlStr, nil)

    // let's add data to the redis db
    gm := gdm.New()
    playerData := dt.NewWithId(PLAYERID) 
    done := gm.StorePlayerData(GAMENAME, playerData)
    if !done {
        t.Errorf("Data shall have been stored\n")
    } 

    // continue
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleUpdatePlayerScore)

 
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    newHiScore, err := gm.GetPlayerHighScore(GAMENAME, PLAYERID)
    if err != nil {
        t.Errorf("TestHandleUpdatePlayerScore Error: GetPlayerHighScore, Error: %v\n", err)
    } else {
        if newHiScore != scoreF {
            t.Errorf("TestHandleUpdatePlayerScore Error: newHiScore=%d. Added score=%d\n", newHiScore, score)
        }
        playerDataStrFromDB, success := gm.GetPlayerData(GAMENAME, PLAYERID)
        if !success {
            t.Errorf("TestHandleUpdatePlayerScore Error: GetPlayerData\n")
        } else {
            playerData.A = scoreF
            playerDataStr :=  dt.Str(playerData)
            if playerDataStrFromDB != playerDataStr {
                t.Errorf("TestHandleUpdatePlayerScore Error: playerDataStr from redis: %s\n\tplayerDataStr shall have been %s\n", playerDataStrFromDB, playerDataStr)
            }
        }
    }

    gm.DelPlayerData(GAMENAME, PLAYERID)
    gm.DeletePlayerScore(GAMENAME, PLAYERID)
}

func TestHandleGetPlayerRank(t *testing.T) {
    var url         bytes.Buffer

    gameName    := GAMENAME
    playerId    := PLAYERID

    url.WriteString("http://gogameserver.com/results?")
    url.WriteString(GAMENAMESHORT + "="+gameName + "&")
    url.WriteString(PLAYERIDSHORT + "="+playerId + "&") 

    urlStr   := url.String()

    req, _ := http.NewRequest("GET", urlStr, nil)

    // let's add data to the redis dbb
    gm := gdm.New()
    scores := []float64{2,1,7,4, 3}
    for i:=0; i<5; i++ {
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
        playerData := dt.NewWithId(PLAYERIDS[i]) 
        playerData.A = scores[i]
        gm.StorePlayerData(GAMENAME, playerData)
        gm.StorePlayerScore(GAMENAME, scores[i], PLAYERIDS[i])
    }

    // continue
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleGetPlayerRank)

    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    rank := gm.GetPlayerRank(GAMENAME, PLAYERIDS[2])
    if rank != 1{
        t.Errorf("TestGetPlayerRank Error: Player rank should have been 1 but is : %d\n", rank)
    }

    for i:=0; i<5; i++ {
        gm.DelPlayerData(GAMENAME, PLAYERIDS[i])
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
    } 
}

func TestHandleGetTopScorersDaily(t *testing.T) {
    gameName    := GAMENAME
    playerId    := PLAYERID

    url.WriteString("http://gogameserver.com/results?")
    url.WriteString(GAMENAMESHORT + "="+gameName + "&")

    
}