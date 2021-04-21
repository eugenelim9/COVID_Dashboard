package sessions

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

//var ctx = context.Background()

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	// Client: client
	// TODO: use param
	return &RedisStore{
		client,
		sessionDuration,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//marshal the `sessionState` to JSON and save it in the redis database,
	//using `sid.getRedisKey()` for the key.
	//return any errors that occur along the way.
	data, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}
	if err := rs.Client.Set(sid.getRedisKey(), data, rs.SessionDuration).Err(); err != nil {
		return err
	}
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter

	//for extra-credit using the Pipeline feature of the redis
	//package to do both the get and the reset of the expiry time
	//in just one network round trip!
	pipe := rs.Client.Pipeline()
	get := pipe.Get(sid.getRedisKey())
	expire := pipe.Expire(sid.getRedisKey(), rs.SessionDuration)
	pipe.Exec()

	data, getErr := get.Bytes()
	if getErr != nil {
		return ErrStateNotFound
	}

	if expire.Err() != nil {
		return expire.Err()
	}

	return json.Unmarshal(data, sessionState)

}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	//delete the data stored in redis for the provided SessionID
	return rs.Client.Del(sid.getRedisKey()).Err()
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
