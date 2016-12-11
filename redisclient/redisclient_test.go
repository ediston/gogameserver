package redisclient_test

import (
    rcl "gogameserver/redisclient" 
    "testing"

    "reflect"
    "time"
)


const tempKeyStr string = "00NeverAddThiskeytemp"
const keyStr string  = "00NeverAddThiskey"
const valStr string  = "00NeverAddThisVal"
const setName string = "00NeverAddThisSet"
const setKey string  = "00NeverAddThisSetKey"
var tempStrs  = [] string{"00NeverAddThiskey0", "00NeverAddThiskey1",  "00NeverAddThiskey2", "00NeverAddThiskey3", "00NeverAddThiskey4",  "00NeverAddThiskey5"}

func TestSaveKeyValTemporary(t *testing.T) {
    rc := rcl.New()
    rc.SaveKeyValTemporary(tempKeyStr, valStr, 1*time.Second) // 10 seconds 10*1000 000 000
    exists,_ := rc.KeyExists(tempKeyStr)
    if !exists {
        t.Errorf("Key should exist!")
    }
    time.Sleep(3 * time.Second)
    exists,_ = rc.KeyExists(tempKeyStr)
    if exists {
        t.Errorf("Key should be deleted!")
        rc.DelKey(tempKeyStr)
    }
}

func TestSaveKeyValForever(t *testing.T) {
    rc := rcl.New()
    rc.SaveKeyValForever(keyStr, valStr)
    exists,_ := rc.KeyExists(keyStr)
    if !exists {
        t.Errorf("Key should exist!")
    }
    rc.DelKey(keyStr)
}

func TestGetVal(t *testing.T) {
    rc := rcl.New()
    rc.SaveKeyValForever(keyStr, valStr)
    tempVal, _ := rc.GetVal(keyStr) 
    if valStr != tempVal{
        t.Errorf("Key should exist and be equal to %s!", valStr)
    }
    rc.DelKey(keyStr)
}

func TestAddToSet(t *testing.T) {
    rc := rcl.New()
    score := 12.0
    rc.AddToSet(setName, score, setKey)
    tempScore, _ := rc.GetScore(setName, setKey)
    if tempScore != score {
        t.Errorf("Stored Score is wrong!")
    }
    rc.RemScore(setName, setKey)
}

func TestGetTop(t *testing.T) {
    rc := rcl.New()
    scores := []float64{2,1,7,4, 3}
    rev_sorted_scores := []float64{7,4,3}
    for i:=0; i<5; i++ {
        rc.AddToSet(setName, scores[i], tempStrs[i])
    }

    top3,_ := rc.GetTop(setName, 3)
    s := reflect.ValueOf(top3)

    for i:=0; i<3; i++ {
        f  := s.Index(i).Field(0)
        if rev_sorted_scores[i] !=  f.Interface() {
            t.Errorf("%d: %s = %v\n", i, f.Type(), f.Interface())
        }
    }
    for i:=0; i<5; i++ {
        rc.RemScore(setName, tempStrs[i])
    }
}

func TestGetRank(t *testing.T) {
    rc := rcl.New()
    scores := []float64{2,1,7,4, 3}
    for i:=0; i<5; i++ {
        rc.AddToSet(setName, scores[i], tempStrs[i])
    }
    rank,_ := rc.GetRank(setName, tempStrs[2])
    if rank != 0{
         t.Errorf("Rank is : %d\n", rank)
    }
    for i:=0; i<5; i++ {
        rc.RemScore(setName, tempStrs[i])
    }
}

func TestGetScore(t *testing.T) {
    rc := rcl.New()
    scores := []float64{2,1,7,4, 3}
    for i:=0; i<5; i++ {
        rc.AddToSet(setName, scores[i], tempStrs[i])
    }
    for i:=0; i<5; i++ {
        score, _ := rc.GetScore(setName, tempStrs[i])
        if score != scores[i] {
            t.Errorf("%d: Expected score: %f. Score is %f\n", i, scores[i], score)
        }
    }
    for i:=0; i<5; i++ {
        rc.RemScore(setName, tempStrs[i])
    }
}

func TestRemScore(t *testing.T) {
    rc := rcl.New()
    rc.AddToSet(setName, 2, tempStrs[0])
    rc.RemScore(setName, tempStrs[0])
    score, _ := rc.GetScore(setName, tempStrs[0])
    if score != 0 {
        t.Errorf("%s exists with score: %f\n", tempStrs[0], score)
    }
}

