package gamedatamanager_test

import (
    dt  "gogameserver/datatypes"
    gdm "gogameserver/gamedatamanager" 
    "testing"
)

const GAMENAME string  = "00dummygame"
const PLAYERID string  = "00playerid"

func  TestStorePlayerData(t *testing.T) {
    gm := gdm.New()
    playerData := dt.NewWithId(PLAYERID) 
    done := gm.StorePlayerData(GAMENAME, playerData)
    if !done {
        t.Errorf("Data shall have been stored\n")
    } 
}

