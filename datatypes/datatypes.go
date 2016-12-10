package datatypes

import (
	"encoding/json"
)

type PlayerData struct {
	I string  	//id   
	IT string	//idType
	N string	// name
	J int64	 `json:",string"`	//joinDate
	L int64	`json:",string"`    //lastPlayed string
	NP int64 `json:",string"`	//numOfTimesPlayed int64
	A float64	 `json:",string"`	//allTimeHighScore int64
	O []string	//oldIds []string
	OT []string	//oldIdtypes []string
	D [] string //devices [] string
	OS []string	//os []string
}

func New(	id string, 
			idType string,
			name string,
			joinDate int64,
			lastPlayed int64,
			numOfTimesPlayed int64,
			allTimeHighScore float64,
			oldIds []string,
			oldIdtypes []string,
			devices [] string,
			os []string	) (PlayerData) {

    return PlayerData{
        I  : id, 
		IT  : idType,
		N  : name,
		J  : joinDate,
		L  : lastPlayed,
		NP  : numOfTimesPlayed,
		A   : allTimeHighScore,
		O  : oldIds,
		OT  : oldIdtypes,
		D   : devices,
		OS  : os,
    }
}

func NewWithId(id string) (PlayerData){
    return PlayerData{
        I  : id, 
		IT  : "",
		N  : "",
		J  : 0,
		L  : 0,
		NP  : 0,
		A   : 0,
		O  : []string{},
		OT  : []string{},
		D   : []string{},
		OS  : []string{},
    }
}

type ResponseData struct {
	PlayerIds []string
	Scores []float64
}

func StrRD (rd ResponseData) string {
	b, _ := json.Marshal(rd)
	return string(b)
}

// ----
func  Str (pd PlayerData) string {
	b, _ := json.Marshal(pd)
	return string(b)
}

func JsonFromStr(s string) PlayerData {
	var ug PlayerData
    err := json.Unmarshal([]byte(s), &ug)
	if err == nil {
		return ug
    } else {
    	return NewWithId("-1")
    }
}
