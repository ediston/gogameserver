package redisclient

import (
    "time"

    "gopkg.in/redis.v5"
)

type redisClient struct {
	client *redis.Client
}

func (rc redisClient) Init() {
	rc.setClient()
}

func (rc redisClient) SetClient(){
	if rc.client != nil{
		return
	}
    rc.client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        DialTimeout:  10 * time.Second,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        PoolSize:     10000,
        PoolTimeout:  30 * time.Second,
    })
}

// dur = int64 nanosecond count.
func (rc redisClient) SaveKeyValTemporary(key string, val interface{}, dur time.Duration) error{
	setClient()
	err := rc.client.Set(key, val, dur).Err()
	if err != nil {
		return error
	}
	return  nil
}

// 
func (rc redisClient) SaveKeyValForever(key string, val interface{}) error{
	setClient()
	return  saveKeyValTemporary(key, val, 0)
}

// 
func (rc redisClient) DelKey(key string) interface{}{
	setClient()
	return  rc.client.Del(key)
}

// 
func (rc redisClient) KeyExists(key string) bool{
	setClient()
	return  rc.client.Exists(key)
}

// 
func (rc redisClient) GetVal(key string) (interface{}, Error){
	setClient()
	val, err := rc.client.Get(key).Result()
	return val, err
}

func (rc redisClient) AddToSet(setName string, value int64, key interface{}){
	setClient()
	val, err := rc.client.ZAdd(setName, value, key)
	return val, err
}

func (rc redisClient) GetTop(setName string, topAmount int64) (interface{}, Error){
	setClient()
	if topAmount <= 0 {
		topAmount = 1
	}
	val, err := rc.client.ZRangeWithScores(setName, start, topAmount-1)
	return val, err
}

func (rc redisClient) GetRank(setName string, key string) (interface{}, Error){
	setClient()
	val, err := rc.client.ZRangeWithScores(setName, key)
	return val, err
}

func (rc redisClient) GetScore(setName string, key string) (interface{}, Error){
	setClient()
	val, err := rc.client.ZScore(setName, key)
	return val, err
}

func (rc redisClient) RemScore(setName string, key string) (interface{}, Error){
	setClient()
	val, err := rc.client.ZRem(setName, key)
	return val, err
}

