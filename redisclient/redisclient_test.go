package redisclient_test

import (
    "string"
    "redisclient"
    "testing"
)


const keyStr string  = "00NeverAddThiskey"
const valStr string  = "00NeverAddThisVal"
const setName string = "00NeverAddThisSet"
const setKey string  = "00NeverAddThisSetKey"
var tempStrs  = [] string{"00NeverAddThiskey0", "00NeverAddThiskey1",  "00NeverAddThiskey2", "00NeverAddThiskey3", "00NeverAddThiskey4",  "00NeverAddThiskey5"}


func TestSaveKeyValTemporary(t *testing.T) {
    var rc redisClient
    ttl := 10
    rc.SaveKeyValTemporary(keyStr, valStr, ttl*Second) // 10 seconds 10*1000 000 000
    if !rc.KeyExists() {
        t.Errorf("Key should exist!")
    }
    Sleep((ttl+1) * Second)
    if rc.KeyExists(keyStr) {
        t.Errorf("Key should be deleted!")
        rc.DelKey(keyStr)
    }
}

// 
func TestSaveKeyValForever(t *testing.T) {
    var rc redisClient
    rc.SaveKeyValForever(keyStr, valStr)
    if !rc.KeyExists(keyStr) {
        t.Errorf("Key should exist!")
    }
    rc.DelKey(keyStr)
}

// 
func TestGetVal(t *testing.T) {
    var rc redisClient
    rc.SaveKeyValForever(keyStr, valStr)
    tempVal, _ := rc.GetVal(keyStr) 
    if valStr != tempVal{
        t.Errorf("Key should exist and be equal to %s!", valStr)
    }
    rc.DelKey(keyStr)
}

func TestAddToSet(t *testing.T) {
    var rc redisClient
    score := 12
    rc.AddToSet(setName, score, setKey)
    tempScore, err := rc.GetScore(setName, setKey)
    if tempScore != score {
        t.Errorf("Stored Score is wrong!")
    }
    rc.RemScore(setName, setKey)
}

func TestGetTop(t *testing.T) {
    var rc redisClient
    var scores := []int64{2,1,7,4,0}
    for i:=0; i<5; i++ {
        rc.AddToSet(setName, scores[i], tempStrs[i])
    }
    top3 := rc.GetTop(setName, 3)
    t.Errorf("%v\n", top3)
    for i:=0; i<5; i++ {
        rc.RemScore(setName, scores[i])
    }
}


