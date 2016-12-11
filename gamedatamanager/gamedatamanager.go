package gamedatamanager


import (
	"log"
	"reflect"
	"encoding/json"
	dt  "gogameserver/datatypes" 
    rcl "gogameserver/redisclient" 
)

type GameManager struct {
	rc *rcl.RedisClient
}

func New() (gm * GameManager) {
    return &GameManager{
        rc: rcl.New(),
    }
}

func (gm * GameManager) GetPlayerData(gameName string, playerId string) (string, bool) {
	playerDataStr, err := gm.rc.GetVal(gameName+playerId)
	if err != nil {
		go log.Printf("\nERROR:GetPlayerData: Game %s, playerId %s not found, err:%v", gameName, playerId, err)
		return "player not found", false
	} else {
		return playerDataStr, true
	}
}

// store player data
func (gm * GameManager) StorePlayerData(gameName string, playerData dt.PlayerData) (bool){
	err := gm.rc.SaveKeyValForever(gameName+playerData.I, dt.Str(playerData))
	if err != nil {
		go log.Printf("\nERROR:StorePlayerData: Game %s, playerData %v, err:%v", gameName, playerData, err)
		return false
	} else {
		go log.Printf("\nInfo: Success StorePlayerData: Game %s, playerData %v", gameName, playerData)
		return true
	}
}

// store player new score for a week
func (gm * GameManager) StorePlayerScoreWeek(gameName string, score float64, playerId string) {
	gm.rc.AddToSet(gameName+"week", score, playerId)
}

// store player new score
func (gm * GameManager) StorePlayerScore(gameName string,  score float64, playerId string) (bool){
	currHiScore, err := gm.GetPlayerHighScore(gameName, playerId)
	if err != nil {
		go log.Printf("ERROR:StorePlayerScore: Game %s, playerId %s not found. err:%v", gameName, playerId, err)
		return false
	} else {
		if currHiScore < score {
			pDataStr, err := gm.GetPlayerData(gameName, playerId)
			if err {
				return false
			}
			pData := dt.JsonFromStr(pDataStr)
			pData.A = score 
			// a go routine to update
			go gm.StorePlayerData(gameName, pData)
			go gm.StorePlayerScoreWeek(gameName, score, playerId)

			gm.rc.AddToSet(gameName, score, playerId)
			go log.Printf("Info:StorePlayerScore: SUCESS currHiScore:%.6f, newScore:%.6f, Game %s, playerId %s", currHiScore, score, gameName, playerId)
		} else {
			go log.Printf("Info:StorePlayerScore: SUCESS currHiScore:%.6f, newScore:%.6f, Game %s, playerId %s", currHiScore, score, gameName, playerId)
		}
		return true
	}
}

// get score
func (gm * GameManager) GetPlayerHighScore(gameName string, playerId string) (float64, error) {
	return gm.rc.GetScore(gameName, playerId)
}

// get player rank
func (gm * GameManager) GetPlayerRank(gameName string, playerId string) (int64) {
	rank, err := gm.rc.GetRank(gameName, playerId)
	if err != nil {
		go log.Printf("Error:GetPlayerRank: Game %s, playerId %s", gameName, playerId)
		return -1;
	} else {
		return rank;
	}
}

// get top x
func (gm * GameManager) GetTopPlayers(gameName string, top int64) (string) {
    var topPlayersWithScores dt.ResponseData

	topResults, err := gm.rc.GetTop(gameName, top)
	if err != nil {
		go log.Printf("Error:GetTopPlayers: Game %s", gameName)
		return ""
	} else {
		go log.Printf("Error:GetTopPlayers: Game %s", gameName)
	}
    topResultsVal := reflect.ValueOf(topResults)
    resultCount := topResultsVal.Len()
    topPlayersWithScores.PlayerIds = make([]string, resultCount)
    topPlayersWithScores.Scores = make([]float64, resultCount)
    for i:=0; i<resultCount; i++ {
    	_score,_  := topResultsVal.Index(i).Field(0).Interface().(float64)
    	_pid,_  := topResultsVal.Index(i).Field(1).Interface().(string)
        topPlayersWithScores.PlayerIds[i] = _pid
        topPlayersWithScores.Scores[i] = _score
    }
    b, jerr := json.Marshal(topPlayersWithScores)
	if jerr == nil {
		return string(b)
	} else {
		return "json error"
	}
}

// get top weekly 1000
func (gm * GameManager) GetTopPlayersWeekly(gameName string) (string) {
	return gm.GetTopPlayers(gameName+"weekly", 1000)
}

// get rank among friends
func (gm * GameManager) GetScoreOfFriends(gameName string, playerId string, friendIds []string) (string) {
	var topPlayersWithScores dt.ResponseData
	playerScore, err := gm.GetPlayerHighScore(gameName, playerId)
	if err != nil {
		go log.Printf("Error:GetRankAmongFriends: Error: %v", err)
		return ""
	}
    
	totalCount := len(friendIds) + 1
    topPlayersWithScores.PlayerIds = make([]string, totalCount)
    topPlayersWithScores.Scores = make([]float64, totalCount)
	for i:=1; i<totalCount; i++ {
		topPlayersWithScores.PlayerIds[i] = friendIds[i-1]
		topPlayersWithScores.Scores[i] = -1
		score, err := gm.GetPlayerHighScore(gameName, friendIds[i-1])
		if err != nil {
			go log.Printf("Error:GetTopPlayers: Game %s", gameName)
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


