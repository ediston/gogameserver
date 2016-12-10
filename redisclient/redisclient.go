package redisclient

import (
    "time"
    "gopkg.in/redis.v5"
)

type RedisClient struct {
	client *redis.Client
}

func New() (rc * RedisClient) {
    return &RedisClient{
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


func (rc *RedisClient) SetClient(){
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
func (rc *RedisClient) SaveKeyValTemporary(key string, val interface{}, dur time.Duration) error{
	rc.SetClient()
	err := rc.client.Set(key, val, dur).Err()
	if err != nil {
		return err
	}
	return  nil
}

// 
func (rc *RedisClient) SaveKeyValForever(key string, val interface{}) error{
	rc.SetClient()
	return  rc.SaveKeyValTemporary(key, val, 0)
}

// 
func (rc *RedisClient) DelKey(key string) (int64, error){
	rc.SetClient()
	return  rc.client.Del(key).Result()
}

// 
func (rc *RedisClient) KeyExists(key string) (bool, error){
	rc.SetClient()
	return  rc.client.Exists(key).Result()
}

// 
func (rc *RedisClient) GetVal(key string) (string, error){
	rc.SetClient()
	return rc.client.Get(key).Result()
}

func (rc *RedisClient) AddToSet(setName string, Score float64, Member interface{}) (int64, error){
	rc.SetClient()
	return rc.client.ZAdd(setName, redis.Z{Score, Member}).Result()
}

//([]Z, error)
/*
type Z struct {
    Score  float64
    Member string
}
type ZRangeByScore struct {
    Min, Max string

    Offset, Count int64
}

*/

// returns ([]Z, error)
func (rc *RedisClient) GetTop(setName string, topAmount int64) (interface{}, error){
	rc.SetClient()
	if topAmount <= 0 {
		topAmount = 1
	}
	return rc.client.ZRevRangeWithScores(setName, 0, topAmount-1).Result()
}

// Rank starts from 0
func (rc *RedisClient) GetRank(setName string, key string) (int64, error){
	rc.SetClient()
	return rc.client.ZRevRank(setName, key).Result()
}

func (rc *RedisClient) GetScore(setName string, key string) (float64, error){
	rc.SetClient()
	return rc.client.ZScore(setName, key).Result()
}

func (rc *RedisClient) RemScore(setName string, key string)  (int64, error){
	rc.SetClient()
	return rc.client.ZRem(setName, key).Result()
}

