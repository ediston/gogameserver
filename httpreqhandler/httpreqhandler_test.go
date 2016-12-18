package httpreqhandler

import (
    dt   "gogameserver/datatypes"
    gdm  "gogameserver/gamedatamanager"
    util "gogameserver/util"

    "testing"
    "strconv"
    "bytes"
    "net/http"
    "net/http/httptest"
    "encoding/json"
    "math/rand"
    "sort"
    "time"
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


    var resp map[string]interface{}

    json.Unmarshal([]byte(rr.Body.String()), &resp)
    rankRec, _ := resp[PLAYERRANK].(float64)

    rank := gm.GetPlayerRank(GAMENAME, PLAYERIDS[2])
    if rank != int64(rankRec){
        t.Errorf("TestGetPlayerRank Error: Player rank should have been 1 but is : %f\n", resp[PLAYERRANK])
    }

    for i:=0; i<5; i++ {
        gm.DelPlayerData(GAMENAME, PLAYERIDS[i])
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
    } 
}

/*
topDailyScores  : [{iXYOT 84.54126154139138} {QrnOZ 82.2631082171108} {EvQHe 77.81027489442664} {eiOJG 70.9182207621995} {XjbBb 68.52574727578084}]
*/

func TestHandleGetTopScorersDaily(t *testing.T) {
    var url         bytes.Buffer
    gameName    := GAMENAME

    topCount := 5
    topAmountStr := strconv.Itoa(topCount)

    url.WriteString("http://gogameserver.com/results?")
    url.WriteString(GAMENAMESHORT + "="+gameName + "&")
    url.WriteString(TYPESHORT + "="+DAY + "&")
    url.WriteString(PLAYERIDSHORT + "=" + "1" + "&")
    url.WriteString(TOPAMOUNTSHORT + "=" + topAmountStr)

    urlStr   := url.String()

    req, _ := http.NewRequest("GET", urlStr, nil)

    // let's add data to the redis dbb
    gm := gdm.New()

    gm.DelKey(GAMENAME+util.GetDate(0))
    gm.DelKey(GAMENAME)


    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    scores := make([] float64, 0)
    names  := make([] string, 10)
    id     := make([] string, 10)

    for i:=0; i<10; i++ {
        names[i] = util.RandStringRunes(5)
        id[i] = strconv.Itoa(i)
        playerData := dt.NewWithId(id[i]) 
        playerData.N = names[i]
        gm.StorePlayerData(GAMENAME, playerData)
    }

    playeScores := make([]dt.PlayerScore, 0)

    // for 1 day
    d := 0
    for i:=0; i<10; i++ {
        score := r.Float64()*100
        scores = append(scores, score)
        playeScores = append(playeScores, dt.PlayerScore{names[i], score})
        gm.StorePlayerScoreOnADay(GAMENAME, score, id[i], d) 
    }
   
    sort.Sort(dt.ByScoreRev(playeScores))

    // continue
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleGetTopScorers)

    handler.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var resp map[string]interface{}

    json.Unmarshal([]byte(rr.Body.String()), &resp)

    topPlayersJsonStr := resp[DAY].(string)
        
    var topDailyScores []dt.PlayerScore
    json.Unmarshal([]byte(topPlayersJsonStr), &topDailyScores)

    for i:=0; i<topCount; i++{
        if playeScores[i].N != topDailyScores[i].N || playeScores[i].S != topDailyScores[i].S {
            t.Errorf("Error: TestHandleGetTopScorersDaily,  playeScores[i] : %v\n, topDailyScores[i]: %v\n", playeScores[i], topDailyScores[i])
        }
    }

    for i:=0; i<10; i++ {
        gm.DelPlayerData(GAMENAME, id[i])
        gm.DeletePlayerScoreOnADay(GAMENAME ,  id[i], d)
    }
    gm.DelKey(GAMENAME+util.GetDate(0))
    gm.DelKey(GAMENAME)
}

/*
topPlayersJsonStr=[{"N":"XqXBF","S":99.2909190190963},
{"N":"ORLng","S":97.12361457169413},{"N":"KQhsx","S":95.7985843790939},
{"N":"FgPKG","S":94.5300503583963},{"N":"iDDSC","S":91.09387827554075}]
*/
func TestHandleGetTopScorersWeekly(t *testing.T) {
    var url  bytes.Buffer
    gameName := GAMENAME

    topCount := 5
    topAmountStr := strconv.Itoa(topCount)

    url.WriteString("http://gogameserver.com/results?")
    url.WriteString(GAMENAMESHORT + "="+gameName + "&")
    url.WriteString(TYPESHORT + "="+WEEK + "&")
    url.WriteString(PLAYERIDSHORT + "=" + "1" + "&")
    url.WriteString(TOPAMOUNTSHORT + "=" + topAmountStr)

    urlStr   := url.String()

    req, _ := http.NewRequest("GET", urlStr, nil)

    // let's add data to the redis dbb
    gm := gdm.New()

    r := rand.New(rand.NewSource(time.Now().UnixNano()))
   
    scores := make([] float64, 0)
    names  := make([] string, 10)
    id     := make([] string, 10)

    for i:=0; i<10; i++ {
        names[i] = util.RandStringRunes(5)
        id[i] = strconv.Itoa(i)
        playerData := dt.NewWithId(id[i]) 
        playerData.N = names[i]
        gm.StorePlayerData(GAMENAME, playerData)
    }


    playerMaxScore := make(map[string] float64)
    for i:=0; i<10; i++ {
        playerMaxScore[names[i]] = 0
    }

    // for 7 days
    for d:=0; d<7; d++ {
        gm.DelKey(GAMENAME+util.GetDate(d))
        for i:=0; i<10; i++ {
            score := r.Float64()*100
            scores = append(scores, score)
            if playerMaxScore[names[i]] < score {
               playerMaxScore[names[i]] = score
            }
            gm.StorePlayerScoreOnADay(GAMENAME, score, id[i], d) 
        }
    } 

    playeScores := make([]dt.PlayerScore, 0)
    for k,v := range playerMaxScore {
        playeScores = append(playeScores, dt.PlayerScore{k, v})
    }
   
    sort.Sort(dt.ByScoreRev(playeScores))


    // continue
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleGetTopScorers)

    handler.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var resp map[string]interface{}

    json.Unmarshal([]byte(rr.Body.String()), &resp)

    topPlayersJsonStr := resp[WEEK].(string)

    var topWeeklyScores []dt.PlayerScore
    json.Unmarshal([]byte(topPlayersJsonStr), &topWeeklyScores)
    

    for i:=0; i<topCount; i++{
        if playeScores[i].N != topWeeklyScores[i].N || playeScores[i].S != topWeeklyScores[i].S {
            t.Errorf("Error: TestHandleGetTopScorersWeekly,  playeScores[i] : %v\n, topWeeklyScores[i]: %v\n", playeScores[i], topWeeklyScores[i])
        }
    }

    for i:=0; i<10; i++ {
        gm.DelPlayerData(GAMENAME, id[i])
    }

    for d:=0; d<7; d++ {
        for i:=0; i<10; i++ {
            gm.DeletePlayerScoreOnADay(GAMENAME ,  id[i], d)
        }
    }
}

