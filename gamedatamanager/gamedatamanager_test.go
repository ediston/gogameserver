package gamedatamanager_test

import (
    dt  "gogameserver/datatypes"
    gdm "gogameserver/gamedatamanager" 
    "testing"
)

const GAMENAME string  = "00dummygame"
const PLAYERID string  = "00playerid"
const PLAYERID2 string  = "002playerid"

func TestStorePlayerData(t *testing.T) {
    gm := gdm.New()
    playerData := dt.NewWithId(PLAYERID) 
    done := gm.StorePlayerData(GAMENAME, playerData)
    if !done {
        t.Errorf("Data shall have been stored\n")
    } 
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
}

func TestStorePlayerScore(t *testing.T) {
    gm := gdm.New()
    playerData := dt.NewWithId(PLAYERID2) 
    playerData.A = 10.0
    newScore := 10000.0
    gm.DeletePlayerScore(GAMENAME, PLAYERID2)
    gm.StorePlayerData(GAMENAME, playerData)
    success := gm.StorePlayerScore(GAMENAME, newScore, PLAYERID2)
    if !success {
        t.Errorf("Error: StorePlayerScore: Data shall have been found\n")
    } else {
        newHiScore, err := gm.GetPlayerHighScore(GAMENAME, PLAYERID2)
        if err != nil {
            t.Errorf("TestStorePlayerScore Error: GetPlayerHighScore, Error: %v\n", err)
        } else {
            if newHiScore != newScore {
                t.Errorf("TestStorePlayerScore Error: newHiScore=%d. Added newScore=%d\n", newHiScore, newScore)
            }
            playerDataStrFromDB, success := gm.GetPlayerData(GAMENAME, PLAYERID2)
            if !success {
                t.Errorf("TestStorePlayerScore Error: GetPlayerData\n", )
            } else {
                playerData.A = newScore
                playerDataStr :=  dt.Str(playerData)
                if playerDataStrFromDB != playerDataStr {
                    t.Errorf("TestStorePlayerScore Error: playerDataStr from redis: %s\n\tplayerDataStr shall have been %s\n", playerDataStrFromDB, playerDataStr)
                }
            }
        }
    }
}

