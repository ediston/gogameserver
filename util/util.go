package util

import (
   "time"
   "net/http"
   "strconv"
   "math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func CurrentDate() string {
    return GetDate(0)
}

func GetDate(daysOld int) string {
    t := time.Now().AddDate(0, 0,-daysOld)
	return t.Format("2006-01-02")
}

func URLParamStr(req *http.Request, paramName string, defParamVal string) string {
	vals, ok := req.URL.Query()[paramName]
	if ok {
		return vals[0]
	}
	return defParamVal
}

func URLParamInt(req *http.Request, paramName string, defParamVal int64) int64 {
	vals, ok := req.URL.Query()[paramName]
	if ok {
		i, _ := strconv.ParseInt(vals[0], 10, 64)
		return i
	}
	return defParamVal
}

func URLParamFloat(req *http.Request, paramName string, defParamVal float64) float64 {
	vals, ok := req.URL.Query()[paramName]
	if ok {
		f, _ := strconv.ParseFloat(vals[0], 64)
		return f
	}
	return defParamVal
}

func RandStringRunes(n int) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}