package gamedatamanager_test

import (
    dt  "gogameserver/datatypes"
    gdm "gogameserver/gamedatamanager" 
    "testing"
)

const GAMENAME string  = "00dummygame"
const PLAYERID string  = "00playerid"

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
        t.Errorf("Data shall have been stored\n")
    } 
    playerDataStr := dt.Str(playerData)
    playerDataStrFromDB, found := gm.GetPlayerData(GAMENAME, PLAYERID)
    if !found {
        t.Errorf("Data shall have been found\n")
    } else {
        if playerDataStrFromDB != playerDataStr {
            t.Errorf("Error: playerDataStr from redis: %s\n\tplayerDataStr shall have been %s\n", playerDataStrFromDB, playerDataStr)
        }
    }
}

