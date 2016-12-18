package gamedatamanager


import (
	"log"
	"reflect"
	"encoding/json"
	"sort"

	"gogameserver/util"
	dt  "gogameserver/datatypes" 
    rcl "gogameserver/redisclient" 
)

const REDIS_NIL string  = "redis: nil"

type GameManager struct {
	rc *rcl.RedisClient
}

func New() (gm * GameManager) {
    return &GameManager{
        rc: rcl.New(),
    }
}

func (gm * GameManager) DelKey(key string) (int64, bool) {
	redisRet, err := gm.rc.DelKey(key)
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("\nERROR:DelKey| Key: %s not found, err:%v", key, err)
		return -1, false
	} else {
		return redisRet, true
	}
}

func (gm * GameManager) GetPlayerData(gameName string, playerId string) (string, bool) {
	gameName = gameName+playerId
	playerDataStr, err := gm.rc.GetVal(gameName)
	if err != nil && err.Error() == REDIS_NIL {
		go log.Printf("\nERROR:GetPlayerData: Game %s, playerId %s not found, err:%v", gameName, playerId, err)
		return "player not found", false
	} else {
		return playerDataStr, true
	}
}

func (gm * GameManager) DelPlayerData(gameName string, playerId string) (int64, bool) {
	gm.DelPlayerName(gameName, playerId)
	gameName = gameName+playerId
	redisRet, err := gm.rc.DelKey(gameName)
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("\nERROR:DelPlayerData: Game %s, playerId %s not found, err:%v", gameName, playerId, err)
		return -1, false
	} else {
		return redisRet, true
	}
}

// store player data
func (gm * GameManager) StorePlayerData(gameName string, playerData dt.PlayerData) (bool){
	gm.StorePlayerName(gameName, playerData.N, playerData.I)
	gameName = gameName+playerData.I
	err := gm.rc.SaveKeyValForever(gameName, dt.Str(playerData))
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("\nERROR:StorePlayerData: Game %s, playerData %v, err:%v", gameName, playerData, err)
		return false
	} else {
		go log.Printf("\nInfo: Success StorePlayerData: Game %s, playerData %v", gameName, playerData)
		return true
	}
}

// store player new score
func (gm * GameManager) StorePlayerScore(gameName string,  score float64, playerId string) (bool){
	currHiScore, err := gm.GetPlayerHighScore(gameName, playerId)
	
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("ERROR:StorePlayerScore: Game %s, playerId %s not found. err:'%s'", gameName, playerId, err.Error() )
		return false
	} else {
		hiScoretoday, _ := gm.GetPlayerScoreOnDay(gameName,  playerId , 0) 
		if hiScoretoday < score {
			gm.StorePlayerScoreDaily(gameName , score , playerId)
		}
		if currHiScore < score {
			pDataStr, found := gm.GetPlayerData(gameName, playerId)
			if !found {
				return false
			}
			pData := dt.JsonFromStr(pDataStr)
			pData.A = score 
			// a go routine to update 
			gm.StorePlayerData(gameName, pData)
			redisRet, redisErr := gm.rc.AddToSet(gameName, score, playerId)
			if redisErr != nil && redisErr.Error() != REDIS_NIL {
				go log.Printf("Error:AddToSet: SUCESS gameName:%s, score:%f, playerId:%s, redisErr:%v", gameName, score, playerId, redisErr)
				return false
			}	
			go log.Printf("Info:StorePlayerScore: SUCESS currHiScore:%.6f, newScore:%.6f, Game %s, playerId %s, retcode:%d", currHiScore, score, gameName, playerId, redisRet)
		} else {
			go log.Printf("Info:StorePlayerScore: Already high currHiScore:%.6f, newScore:%.6f, Game %s, playerId %s", currHiScore, score, gameName, playerId)
		}
		return true
	}
}

func getKeyForName(gameName string, playerId string)  string {
	return gameName+playerId+"Name"
}

// store player new score
func (gm * GameManager) StorePlayerName(gameName string,  playerName string, playerId string) (bool){
	key := getKeyForName(gameName, playerId)
	err := gm.rc.SaveKeyValForever(key, playerName)
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("\nERROR:StorePlayerName: Game %s, playerName %s, err:%v", gameName, playerName, err)
		return false
	} else {
		go log.Printf("\nInfo: Success StorePlayerName: Game %s, playerName %s", gameName, playerName)
		return true
	}
}

// store player new score
func (gm * GameManager) GetPlayerName(gameName string, playerId string) (string, bool){
	key := getKeyForName(gameName, playerId)
	playerName, err := gm.rc.GetVal(key)
	if err != nil && err.Error() == REDIS_NIL {
		go log.Printf("\nERROR:GetPlayerName: Game %s, playerId %s not found, err:%v", gameName, playerId, err)
		return "playerName not found", false
	} else {
		return playerName, true
	}
}

// store player new score
func (gm * GameManager) DelPlayerName(gameName string,   playerId string) (int64, bool){
	key := getKeyForName(gameName, playerId)
	redisRet, err := gm.rc.DelKey(key)
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("\nERROR:DelPlayerName: Game %s, playerId %s not found, err:%v", gameName, playerId, err)
		return -1, false
	} else {
		return redisRet, true
	}
}

// delete player score
func (gm * GameManager) DeletePlayerScore(gameName string,  playerId string) (bool){
	redisRet, redisErr := gm.rc.RemScore(gameName, playerId)
	if redisErr != nil && redisErr.Error() != REDIS_NIL {
		go log.Printf("Error:DeletePlayerScore: SUCESS gameName:%s, playerId:%s, redisErr:%v", gameName, playerId, redisErr)
		return false
	} else {
		go log.Printf("Info :DeletePlayerScore: SUCESS gameName:%s,  playerId:%s, redisRet:%d", gameName, playerId, redisRet)
	}
	return true
}

// get player rank
func (gm * GameManager) GetPlayerRank(gameName string, playerId string) (int64) {
	rank, err := gm.rc.GetRank(gameName, playerId)
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("Error:GetPlayerRank: Game %s, playerId %s", gameName, playerId)
		return -1;
	} else {
		return rank+1;
	}
}

// get top x
func (gm * GameManager) GetTopPlayers(gameName string, topCount int64) (string) {
	playeScores := make([]dt.PlayerScore, 0)

	topResults, err := gm.rc.GetTop(gameName, topCount)
	if err != nil && err.Error() != REDIS_NIL{
		go log.Printf("Error:GetTopPlayers: Game %s", gameName)
		return "json error"
	}
    topResultsVal := reflect.ValueOf(topResults)
    resultCount := topResultsVal.Len()
    for i:=0; i<resultCount; i++ {
    	_score,_  := topResultsVal.Index(i).Field(0).Interface().(float64)
    	_pid,_  := topResultsVal.Index(i).Field(1).Interface().(string)
    	_name,_ := gm.GetPlayerName(gameName, _pid)
    	playeScores = append(playeScores, dt.PlayerScore{_name, _score})
    }

	sort.Sort(dt.ByScoreRev(playeScores))
 
	if topCount < int64(len(playeScores)) {
		playeScores = playeScores[:topCount]
	}
	go log.Printf("Info: GetTopPlayers: top %d playeScores: %v", topCount, playeScores)
	b, jerr := json.Marshal(playeScores)
	if jerr == nil {
		return string(b)
	} else {
		return "json error"
	}
}

// get score
func (gm * GameManager) GetPlayerHighScore(gameName string, playerId string) (float64, error) {
	return gm.rc.GetScore(gameName, playerId)
}

// store player new score for a week
func (gm * GameManager) StorePlayerScoreOnADay(gameName string, score float64, playerId string, numOfDaysOld int) {
	gm.rc.AddToSet(gameName+util.GetDate(numOfDaysOld), score, playerId)
}

func (gm * GameManager) DeletePlayerScoreOnADay(gameName string,  playerId string, numOfDaysOld int) (bool){
	gameName = gameName+util.GetDate(numOfDaysOld)
	redisRet, redisErr := gm.rc.RemScore(gameName, playerId)
	if redisErr != nil && redisErr.Error() != REDIS_NIL {
		go log.Printf("Error:DeletePlayerScoreOnADay: SUCESS gameName:%s, playerId:%s, redisErr:%v", gameName, playerId, redisErr)
		return false
	} else {
		go log.Printf("Info :DeletePlayerScoreOnADay: SUCESS gameName:%s,  playerId:%s, redisRet:%d", gameName, playerId, redisRet)
	}
	return true
}

// store player new score for a week
func (gm * GameManager) StorePlayerScoreDaily(gameName string, score float64, playerId string) {
	gm.rc.AddToSet(gameName+util.CurrentDate(), score, playerId)
}

// store player new score for a week
func (gm * GameManager) GetPlayerScoreOnDay(gameName string, playerId string, numOfDaysOld int) (float64, error) {
	return gm.GetPlayerHighScore(gameName+util.GetDate(numOfDaysOld), playerId)
}

// get top weekly 1000
func (gm * GameManager) GetTopPlayersOnDay(gameName string, topCount int64, numOfDaysOld int) (string) {
	if numOfDaysOld > 6 {
	  return "";
	}

	playeScores := make([]dt.PlayerScore, 0)

	dateGameName := gameName+util.GetDate(numOfDaysOld)

	topResults, err := gm.rc.GetTop(dateGameName, topCount)
	if err != nil && err.Error() != REDIS_NIL{
		go log.Printf("Error:GetTopPlayersOnDay: dateGameName %s", dateGameName)
		return "json error"
	}
    topResultsVal := reflect.ValueOf(topResults)
    resultCount := topResultsVal.Len()
    for i:=0; i<resultCount; i++ {
    	_score,_  := topResultsVal.Index(i).Field(0).Interface().(float64)
    	_pid,_  := topResultsVal.Index(i).Field(1).Interface().(string)
    	_name,_ := gm.GetPlayerName(gameName, _pid)
    	playeScores = append(playeScores, dt.PlayerScore{_name, _score})
    }

	sort.Sort(dt.ByScoreRev(playeScores))
 
	if topCount < int64(len(playeScores)) {
		playeScores = playeScores[:topCount]
	}
	go log.Printf("Info: GetTopPlayersOnDay: top %d playeScores: %v", topCount, playeScores)
	b, jerr := json.Marshal(playeScores)
	if jerr == nil {
		return string(b)
	} else {
		return "json error"
	}

}

func (gm * GameManager) GetTopPlayersThisWeek(gameName string, topCount int64) (string) {
	playeScores := make([]dt.PlayerScore, 0)
    donePersons := make(map[string] bool) 

	for i:=6; i>=0; i--{
		topResults, err := gm.rc.GetTop(gameName+util.GetDate(i), topCount)
		if err != nil && err.Error() != REDIS_NIL{
			go log.Printf("Error:GetTopPlayersThisWeek: Game %s", gameName)
			continue
		}
	    topResultsVal := reflect.ValueOf(topResults)
	    resultCount := topResultsVal.Len()
	    for i:=0; i<resultCount; i++ {
	    	_score,_  := topResultsVal.Index(i).Field(0).Interface().(float64)
	    	_pid,_  := topResultsVal.Index(i).Field(1).Interface().(string)
	    	_name,_ := gm.GetPlayerName(gameName, _pid)
	    	playeScores = append(playeScores, dt.PlayerScore{_name, _score})
	    }
	}

	sort.Sort(dt.ByScoreRev(playeScores))

	for i:=0; i<len(playeScores) && i<int(topCount); i++ {
		if _,ok := donePersons[playeScores[i].N]; !ok {
			donePersons[playeScores[i].N] = true
		} else {
			playeScores = append(playeScores[:i], playeScores[i+1:]...) // perfectly fine if i is the last element https://play.golang.org/p/hcUEguyiTC
			i--
		}
	}
 
	if topCount < int64(len(playeScores)) {
		playeScores = playeScores[:topCount]
	}
	go log.Printf("Info: GetTopPlayersThisWeek: top %d playeScores: %v", topCount, playeScores)
	b, jerr := json.Marshal(playeScores)
	if jerr == nil {
		return string(b)
	} else {
		return "json error"
	}
}

// get rank among friends
func (gm * GameManager) GetScoreOfFriends(gameName string, playerId string, friendIds []string) (string) {
	var topPlayersWithScores dt.ResponseData
	playerScore, err := gm.GetPlayerHighScore(gameName, playerId)
	if err != nil && err.Error() != REDIS_NIL {
		go log.Printf("Error:GetScoreOfFriends: Error: %v", err)
		return ""
	}
    
	totalCount := len(friendIds) + 1
    topPlayersWithScores.PlayerIds = make([]string, totalCount)
    topPlayersWithScores.Scores = make([]float64, totalCount)
	for i:=1; i<totalCount; i++ {
		topPlayersWithScores.PlayerIds[i] = friendIds[i-1]
		topPlayersWithScores.Scores[i] = -1
		score, err1 := gm.GetPlayerHighScore(gameName, friendIds[i-1])
		if err1 != nil && err1.Error() != REDIS_NIL {
			go log.Printf("Error:GetScoreOfFriends: Game %s, %v", gameName, err1)
		} else {
			topPlayersWithScores.Scores[i] = score
		}
	}
	
	topPlayersWithScores.PlayerIds[0] = playerId
	topPlayersWithScores.Scores[0] =playerScore
	b, jerr := json.Marshal(topPlayersWithScores)
	if jerr == nil {
		return string(b)
	} else {
		return "json error"
	}
}
