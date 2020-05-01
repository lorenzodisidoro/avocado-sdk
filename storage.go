package sdk

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	bolt "go.etcd.io/bbolt"
)

// RedisStorage defines redis client settings
type RedisStorage struct {
	Address  string
	Password string
	DB       int
	client   *redis.Client
}

// BboltStorage defines bbolt client settings
type BboltStorage struct {
	SotoragePath string
	BucketName   string
	client       struct {
		options *bolt.Options
		mode    os.FileMode
	}
}

// StorageClient defines key value storage client type
type StorageClient struct {
	Bbolt *BboltStorage
	Redis *RedisStorage
}

// New create a new instance of storage client
func (sc *StorageClient) New() error {

	if sc.Bbolt != nil && sc.Redis != nil {
		return errors.New("Is not possible instantiate more than one storage client")
	}

	if sc.Redis != nil {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     sc.Redis.Address,  // localhost:6379
			Password: sc.Redis.Password, // no password set
			DB:       sc.Redis.DB,       // use default DB
		})

		sc.Redis.client = redisClient
	}

	if sc.Bbolt != nil {
		options := &bolt.Options{Timeout: 10 * time.Second}
		var mode os.FileMode = 0770

		sc.Bbolt.client.options = options
		sc.Bbolt.client.mode = mode
	}

	return nil
}

// Get all keys
func (sc *StorageClient) getAll() ([][]byte, error) {
	var keys [][]byte
	var err error

	if sc.Bbolt != nil {
		bboltClient, err := bolt.Open(sc.Bbolt.SotoragePath, sc.Bbolt.client.mode, sc.Bbolt.client.options)
		if err != nil {
			return nil, err
		}

		bboltClient.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(sc.Bbolt.BucketName))
			if b != nil {
				c := b.Cursor()
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					keys = append(keys, k)
				}
			}

			return nil
		})

		bboltClient.Close()
	}

	if sc.Redis != nil && sc.Redis.client != nil {
		retrnedInterface, err := sc.Redis.client.Do("KEYS", "*").Result()
		if err != nil {
			return keys, err
		}

		// parse to bytes
		results := retrnedInterface.([]interface{})
		for i := 0; i < len(results); i++ {
			key := results[i].(string)
			keys = append(keys, []byte(key))
		}
	}

	return keys, err
}

// Get return a value by key
func (sc *StorageClient) get(key []byte) ([]byte, error) {
	var value []byte
	var err error

	if sc.Bbolt != nil {
		bboltClient, err := bolt.Open(sc.Bbolt.SotoragePath, sc.Bbolt.client.mode, sc.Bbolt.client.options)
		if err != nil {
			return nil, err
		}

		bboltClient.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(sc.Bbolt.BucketName))
			if bucket == nil {
				err = fmt.Errorf("Bucket %q not found", sc.Bbolt.BucketName)
				return err
			}

			value = bucket.Get(key)

			return nil
		})

		bboltClient.Close()
	}

	if sc.Redis != nil && sc.Redis.client != nil {
		valueString, err := sc.Redis.client.Get(string(key)).Result()
		if err != nil {
			return nil, err
		}

		value = []byte(valueString)
	}

	return value, err
}

// Set put new key value pair
func (sc *StorageClient) set(key, value []byte) error {
	if sc.Bbolt != nil {
		bboltClient, err := bolt.Open(sc.Bbolt.SotoragePath, sc.Bbolt.client.mode, sc.Bbolt.client.options)
		if err != nil {
			return err
		}

		err = bboltClient.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(sc.Bbolt.BucketName))
			if err != nil {
				return err
			}

			err = bucket.Put(key, value)
			if err != nil {
				return err
			}
			return nil
		})

		bboltClient.Close()

		return err
	}

	if sc.Redis != nil && sc.Redis.client != nil {
		err := sc.Redis.client.Set(string(key), string(value), 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
