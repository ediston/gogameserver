package redisclient

import (
    "time"

    "gopkg.in/redis.v5"
)

type redisClient struct {
	client *redis.Client
}

func New() (rc * redisClient) {
    return &redisClient{
        client: redis.NewClient(&redis.Options{
	        Addr:     "localhost:6379",
	        DialTimeout:  10 * time.Second,
	        ReadTimeout:  30 * time.Second,
	        WriteTimeout: 30 * time.Second,
	        PoolSize:     10000,
	        PoolTimeout:  30 * time.Second,
    	}),
    }
}


func (rc *redisClient) SetClient(){
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
func (rc *redisClient) SaveKeyValTemporary(key string, val interface{}, dur time.Duration) error{
	rc.SetClient()
	err := rc.client.Set(key, val, dur).Err()
	if err != nil {
		return err
	}
	return  nil
}

// 
func (rc *redisClient) SaveKeyValForever(key string, val interface{}) error{
	rc.SetClient()
	return  rc.SaveKeyValTemporary(key, val, 0)
}

// 
func (rc *redisClient) DelKey(key string) (int64, error){
	rc.SetClient()
	return  rc.client.Del(key).Result()
}

// 
func (rc *redisClient) KeyExists(key string) (bool, error){
	rc.SetClient()
	return  rc.client.Exists(key).Result()
}

// 
func (rc *redisClient) GetVal(key string) (string, error){
	rc.SetClient()
	return rc.client.Get(key).Result()
}

func (rc *redisClient) AddToSet(setName string, Score float64, Member interface{}) (int64, error){
	rc.SetClient()
	return rc.client.ZAdd(setName, redis.Z{Score, Member}).Result()
}

//([]Z, error)
/*
type Z struct {
    Score  float64
    Member string
}
type ZRangeByScore
type ZRangeByScore struct {
    Min, Max string

    Offset, Count int64
}

*/
func (rc *redisClient) GetTop(setName string, topAmount int64) (interface{}, error){
	rc.SetClient()
	if topAmount <= 0 {
		topAmount = 1
	}
	return rc.client.ZRangeWithScores(setName, 0, topAmount-1).Result()
}

// returns ([]Z, error)
func (rc *redisClient) GetRank(setName string, key string) (int64, error){
	rc.SetClient()
	return rc.client.ZRank(setName, key).Result()
}

func (rc *redisClient) GetScore(setName string, key string) float64{
	rc.SetClient()
	return rc.client.ZScore(setName, key).Val()
}

func (rc *redisClient) RemScore(setName string, key string)  (int64, error){
	rc.SetClient()
	return rc.client.ZRem(setName, key).Result()
}

