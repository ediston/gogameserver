package gamedatamanager_test

import (
    dt   "gogameserver/datatypes"
    gdm  "gogameserver/gamedatamanager"
    util "gogameserver/util"

    "testing"
    "math/rand"
    "time"
    "strconv"
    "encoding/json"

    "sort"
)

const GAMENAME string  = "00dummygame" 
const PLAYERID string  = "00playerid"
var PLAYERIDS  = [] string{"00NeverAddThispid0", "00NeverAddThispid1",  "00NeverAddThispid3", "00NeverAddThispid4", "00NeverAddThispid5"}

func TestStorePlayerData(t *testing.T) {
    gm := gdm.New()
    playerData := dt.NewWithId(PLAYERID) 
    done := gm.StorePlayerData(GAMENAME, playerData)
    if !done {
        t.Errorf("Data shall have been stored\n")
    } 
    gm.DelPlayerData(GAMENAME, PLAYERID)
}

func TestGetPlayerData(t *testing.T) {
    gm := gdm.New()
    playerData := dt.NewWithId(PLAYERID) 
    done := gm.StorePlayerData(GAMENAME, playerData)
    if !done {
        t.Errorf("Error: Data shall have been stored\n")
    } 
    playerDataStr := dt.Str(playerData)
    playerDataStrFromDB, found := gm.GetPlayerData(GAMENAME, PLAYERID)
    if !found {
        t.Errorf("Error: GetPlayerData: Data shall have been found\n")
    } else {
        if playerDataStrFromDB != playerDataStr {
            t.Errorf("Error: playerDataStr from redis: %s\n\tplayerDataStr shall have been %s\n", playerDataStrFromDB, playerDataStr)
        }
    }
    gm.DelPlayerData(GAMENAME, PLAYERID)
}

func TestStorePlayerScore(t *testing.T) {
    gm := gdm.New()
    playerData := dt.NewWithId(PLAYERID) 
    playerData.A = 10.0
    newScore := 10000.0
    gm.DeletePlayerScore(GAMENAME, PLAYERID)
    gm.StorePlayerData(GAMENAME, playerData)
    success := gm.StorePlayerScore(GAMENAME, newScore, PLAYERID)
    if !success {
        t.Errorf("Error: StorePlayerScore: Data shall have been found\n")
    } else {
        newHiScore, err := gm.GetPlayerHighScore(GAMENAME, PLAYERID)
        if err != nil {
            t.Errorf("TestStorePlayerScore Error: GetPlayerHighScore, Error: %v\n", err)
        } else {
            if newHiScore != newScore {
                t.Errorf("TestStorePlayerScore Error: newHiScore=%d. Added newScore=%d\n", newHiScore, newScore)
            }
            playerDataStrFromDB, success := gm.GetPlayerData(GAMENAME, PLAYERID)
            if !success {
                t.Errorf("TestStorePlayerScore Error: GetPlayerData\n")
            } else {
                playerData.A = newScore
                playerDataStr :=  dt.Str(playerData)
                if playerDataStrFromDB != playerDataStr {
                    t.Errorf("TestStorePlayerScore Error: playerDataStr from redis: %s\n\tplayerDataStr shall have been %s\n", playerDataStrFromDB, playerDataStr)
                }
            }
        }
    }
    gm.DelPlayerData(GAMENAME, PLAYERID)
    gm.DeletePlayerScore(GAMENAME, PLAYERID)
}

func TestGetPlayerRank(t *testing.T) {
    gm := gdm.New()
    scores := []float64{2,1,7,4, 3}
    for i:=0; i<5; i++ {
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
        playerData := dt.NewWithId(PLAYERIDS[i]) 
        playerData.A = scores[i]
        gm.StorePlayerData(GAMENAME, playerData)
        gm.StorePlayerScore(GAMENAME, scores[i], PLAYERIDS[i])
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


func TestGetScoreOfFriends(t *testing.T) {
    gm := gdm.New()
    scores := []float64{2,1,7,4, 3}
    for i:=0; i<5; i++ {
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
        playerData := dt.NewWithId(PLAYERIDS[i]) 
        playerData.A = scores[i]
        gm.StorePlayerData(GAMENAME, playerData)
        gm.StorePlayerScore(GAMENAME, scores[i], PLAYERIDS[i])
    }
    scoresOfFriends := gm.GetScoreOfFriends(GAMENAME, PLAYERIDS[0], PLAYERIDS[2:5])
    scoresOfFriendsExpectedStr := "{\"PlayerIds\":[\"00NeverAddThispid0\",\"00NeverAddThispid3\",\"00NeverAddThispid4\",\"00NeverAddThispid5\"],\"Scores\":[2,7,4,3]}"
    if scoresOfFriends != scoresOfFriendsExpectedStr{
        t.Errorf("TestGetScoreOfFriends Error: scoresOfFriends str is : %s\n", scoresOfFriends)
    }

    for i:=0; i<5; i++ {
        gm.DelPlayerData(GAMENAME, PLAYERIDS[i])
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
    }
}

func TestGetTopPlayers(t *testing.T) {
    gm := gdm.New()
    scores := []float64{2,1,7,4, 3}
    names  := make([] string, 5)

    playeScores := make([]dt.PlayerScore, 0)
    gm.DelKey(GAMENAME)

    for i:=0; i<5; i++ {
        names[i] = util.RandStringRunes(5)
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
        playerData := dt.NewWithId(PLAYERIDS[i]) 
        playerData.A = scores[i]
        playerData.N = names[i]
        gm.StorePlayerData(GAMENAME, playerData)
        gm.StorePlayerScore(GAMENAME, scores[i], PLAYERIDS[i])

        playeScores = append(playeScores, dt.PlayerScore{names[i], scores[i]})
    }

    sort.Sort(dt.ByScoreRev(playeScores))
    topCount := 3
    top3 := gm.GetTopPlayers(GAMENAME, int64(topCount) )
    
    var top3Scorers []dt.PlayerScore
    json.Unmarshal([]byte(top3), &top3Scorers)

    for i:=0; i<topCount; i++{
        if playeScores[i].N != top3Scorers[i].N || playeScores[i].S != top3Scorers[i].S {
            t.Errorf("Error: TestGetTopPlayers,  playeScores[i] : %v\n, top3Scorers[i]: %v\n", playeScores[i], top3Scorers[i])
        }
    }

    for i:=0; i<5; i++ {
        gm.DelPlayerData(GAMENAME, PLAYERIDS[i])
        gm.DeletePlayerScore(GAMENAME, PLAYERIDS[i])
    }

} 

func TestGetTopPlayersThisWeek(t *testing.T) {
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

    topCount := 5
    topWeeklyScorersJson := gm.GetTopPlayersThisWeek(GAMENAME, int64(topCount) )
    
    var topWeeklyScorers []dt.PlayerScore
    json.Unmarshal([]byte(topWeeklyScorersJson), &topWeeklyScorers)

    for i:=0; i<topCount; i++{
        if playeScores[i].N != topWeeklyScorers[i].N || playeScores[i].S != topWeeklyScorers[i].S {
            t.Errorf("Error: TestGetTopPlayersThisWeek,  playeScores[i] : %v\n, topWeeklyScorers[i]: %v\n", playeScores[i], topWeeklyScorers[i])
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

