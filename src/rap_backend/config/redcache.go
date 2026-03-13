package config

import (
	"time"

	"github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

type RedisCfg struct {
	MaxIdle     int    `json:"maxidle"`
	DbIndex     int    `json:"dbindex"`
	Address     string `json:"address"`
	MaxActive   int    `json:"maxactive"`
	IdleTimeout int    `json:"idletimeout"`
}

type RedisCache struct {
	rediscfg  RedisCfg
	redclient *redis.Pool
}

func (this *RedisCache) cfgEqual(cfg RedisCfg) bool {
	if this.redclient == nil {
		return false
	}

	return cfg == this.rediscfg
}

func (this *RedisCache) redisInit(cfg RedisCfg) int {
	address := cfg.Address
	dbindex := cfg.DbIndex

	if this.cfgEqual(cfg) {
		seelog.Error("new cfg is equal old cfg")
		return 0
	}

	IdleTimeout := time.Duration(cfg.IdleTimeout)
	dial := func() (redis.Conn, error) {
		connection, err := redis.Dial("tcp", address)
		if err != nil {
			seelog.Error("Dial redis ", address, err)
			return nil, err
		}

		if dbindex != 0 {
			_, err = connection.Do("SELECT", dbindex)
			if err != nil {
				seelog.Error("select redis db index ", dbindex, err)
				return nil, err
			}
		}

		return connection, nil
	}

	this.rediscfg = cfg
	if this.redclient != nil {
		redclient := this.redclient
		go redclient.Close()
	}

	this.redclient = &redis.Pool{
		Dial:        dial,
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: IdleTimeout * time.Second,
	}

	seelog.Infof("redis %s db %d ", address, dbindex)
	return 0
}

func (this *RedisCache) Hdel(key string, argvs []interface{}) (r int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	r = -1
	return this.hdel(key, argvs)
}

func (this *RedisCache) hdel(key string, argvs []interface{}) int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	argvs = append([]interface{}{key}, argvs...)
	_, err := conn.Do("HDEL", argvs...)
	if err != nil {
		seelog.Errorf("redis hdel %s error %s", key, err)
		return -1
	}

	return 0
}

func (this *RedisCache) Hset(key string, field string, value string) int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	seelog.Debugf("hset:%s, %s, %s", key, field, value)

	//argvs = append([]interface{}{key}, argvs...)
	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		seelog.Errorf("redis hset %s error %s", key, err)
		return -1
	}

	return 0
}

func (this *RedisCache) Hmset(key string, argvs []interface{}) (r int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	r = -1
	return this.hmset(key, argvs)
}

func (this *RedisCache) hmset(key string, argvs []interface{}) int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	argvs = append([]interface{}{key}, argvs...)
	_, err := conn.Do("HMSET", argvs...)
	if err != nil {
		seelog.Errorf("redis hmset %s error %s", key, err)
		return -1
	}

	return 0
}

func (this *RedisCache) Expire(key string, maxttl int64) (r int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	r = -1
	return this.expire(key, maxttl)
}

func (this *RedisCache) expire(key string, maxttl int64) int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	_, err := conn.Do("EXPIRE", key, maxttl)
	if err != nil {
		seelog.Errorf("redis expire %s error %s", key, err)
		return -1
	}

	return 0
}

func (this *RedisCache) Hgetall(key string) map[string]string {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	return this.hgetall(key)
}

func (this *RedisCache) hgetall(key string) map[string]string {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return nil
	}

	data, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		seelog.Errorf("redis hgetall %s error %s", key, err)
		return nil
	}

	return data
}

func (this *RedisCache) Ihgetall(key string) map[string]int {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	return this.ihgetall(key)
}

func (this *RedisCache) ihgetall(key string) map[string]int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return nil
	}

	data, err := redis.IntMap(conn.Do("HGETALL", key))
	if err != nil {
		seelog.Errorf("redis hgetall %s error %s", key, err)
		return nil
	}

	return data
}

func (this *RedisCache) Hincrby(key, field string, increment int) (r int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	r = -1
	return this.hincrby(key, field, increment)
}

func (this *RedisCache) hincrby(key, field string, increment int) (r int) {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	ret, err := redis.Int(conn.Do("HINCRBY", key, field, increment))
	if err != nil {
		seelog.Errorf("redis hincrby %s error %s", key, err)
		return -1
	}

	return ret
}

func (this *RedisCache) Rpop(key string) string {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	return this.rpop(key)
}

func (this *RedisCache) rpop(key string, argvs ...interface{}) string {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return ""
	}

	ret, err := redis.String(conn.Do("RPOP", key))
	if err != nil {
		return ""
	}

	return ret
}

func (this *RedisCache) Lpush(key string, argvs ...interface{}) (r int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	r = -1
	return this.lpush(key, argvs...)
}

func (this *RedisCache) lpush(key string, argvs ...interface{}) int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	argvs = append([]interface{}{key}, argvs...)
	_, err := conn.Do("LPUSH", argvs...)
	if err != nil {
		seelog.Errorf("LPUSH %s %s", key, err)
		return -1
	}

	return 0
}

func (this *RedisCache) Rpush(key string, argvs ...interface{}) (r int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	r = -1
	return this.rpush(key, argvs...)
}

func (this *RedisCache) rpush(key string, argvs ...interface{}) int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	argvs = append([]interface{}{key}, argvs...)
	_, err := conn.Do("RPUSH", argvs...)
	if err != nil {
		seelog.Errorf("RPUSH %s %s", key, err)
		return -1
	}

	return 0
}

func (this *RedisCache) Llen(key string) int {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return 0
	}

	ret, err := redis.Int(conn.Do("LLEN", key))
	if err != nil {
		seelog.Errorf("LLEN %s %s", key, err)
		return 0
	}

	return ret
}

func (this *RedisCache) Exists(key string) int {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	ret, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		seelog.Errorf("redis exists %s error %s", key, err)
		return -1
	}

	return ret
}
func (this *RedisCache) Delete(key string) (r int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	r = -1
	return this.delete(key)
}

func (this *RedisCache) delete(key string) int {
	conn := this.redclient.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		seelog.Errorf("get redis connection error %s %s", key, err)
		return -1
	}

	ret, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		seelog.Errorf("redis delete %s error %s", key, err)
		return -1
	}

	return ret
}
