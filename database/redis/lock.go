package redis


import (
	"bytes"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/niuhuan/mirai-bot/utils"
	"runtime"
	"strconv"
	"time"
)

var (
	ErrTimeout = errors.New("lock: obtain timeout")
	BootId     = utils.GetSnowflakeIdString()
)

func GoID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func ContextId() string {
	return fmt.Sprintf("%s::%d", BootId, GoID())
}

type Lock struct {
	Key    string
	Expire time.Duration
	RdPool *redis.Pool
}

const lockScript = `
if (redis.call('exists', KEYS[1]) == 0) then 
  redis.call('hincrby', KEYS[1], ARGV[2], 1); 
  redis.call('pexpire', KEYS[1], ARGV[1]); 
  return nil; 
end; 
if (redis.call('hexists', KEYS[1], ARGV[2]) == 1) then 
  redis.call('hincrby', KEYS[1], ARGV[2], 1); 
  redis.call('pexpire', KEYS[1], ARGV[1]); 
  return nil; 
end; 
return redis.call('pttl', KEYS[1]);
`

const unlockScript = `
if (redis.call('hexists', KEYS[1], ARGV[2]) == 0) then 
  return nil;
end; 
local counter = redis.call('hincrby', KEYS[1], ARGV[2], -1); 
if (counter > 0) then 
  redis.call('pexpire', KEYS[1], ARGV[1]); 
  return 0; 
else 
  redis.call('del', KEYS[1]); 
  return 1; 
end; 
return nil;
`

var lockScriptRedis = redis.NewScript(1, lockScript)
var unlockScriptRedis = redis.NewScript(1, unlockScript)

func (lock *Lock) Unlock() (bool, error) {
	// connect
	conn := lock.RdPool.Get()
	defer conn.Close()
	//
	bool, err := redis.Bool(unlockScriptRedis.Do(conn, lock.Key, lock.Expire.Milliseconds(), ContextId()))
	// already lease when if err == redis.ErrNil
	// reentry counter > 0 when bool is false
	if err == nil && bool {
		conn.Do("PUBLISH", lock.Key, lock.Key)
	}
	return bool, err
}

func TryLock(key string, wait time.Duration, lease time.Duration) (*Lock, error) {
	current := time.Now()
	conn := RdPool.Get()
	defer conn.Close()
	// loop
	for true {
		// lock or reentry
		ttl, err := redis.Uint64(lockScriptRedis.Do(conn, key, lease.Milliseconds(), ContextId()))
		if err == redis.ErrNil {
			return &Lock{
				Key:    key,
				Expire: lease,
				RdPool: RdPool,
			}, nil
		} else if err != nil {
			return nil, err
		}
		// timeout
		currentOff := time.Now()
		if currentOff.Sub(current) >= wait {
			return nil, ErrTimeout
		}
		// subscribe
		func() {
			waitDuration := current.Add(wait).Sub(currentOff)
			ttlDuration := time.Millisecond * time.Duration(ttl)
			if ttlDuration < waitDuration {
				waitDuration = ttlDuration
			}
			subConn := redis.PubSubConn{Conn: RdPool.Get()}
			defer subConn.Close()
			subConn.Subscribe(key)
			for true {
				switch subConn.ReceiveWithTimeout(waitDuration).(type) {
				case redis.Message: // message
				case error: // timeout (or redis error)
					return
				case redis.Subscription: // continue for
				default: // else continue
					continue
				}
			}
		}()
	}
	return nil, ErrTimeout
}
