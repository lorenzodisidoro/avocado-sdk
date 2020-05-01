[![Go Report Card](https://goreportcard.com/badge/github.com/lorenzodisidoro/avocado-sdk)](https://goreportcard.com/report/github.com/lorenzodisidoro/avocado-sdk)
[![Build Status](https://travis-ci.com/lorenzodisidoro/avocado-sdk.svg?branch=master)](https://travis-ci.com/lorenzodisidoro/avocado-sdk)

# Avocado SDK
Avocado SDK encrypt with RSA keys a values ​​to be store in the key-value store.

# Supported
- [Redis](https://redis.io/)
- [BBolt](https://github.com/boltdb/bolt)

# How to
## Install
To start using Avocado SDK, install Go and run go get command as a follow
```sh
go get github.com/lorenzodisidoro/avocado-sdk
```

## Use
To use SDK methods import in your GO file
```go
import sdk "github.com/lorenzodisidoro/avocado-sdk"
```

### Storage configuration
Create the storage configuration

#### BBolt
```go
storage := &sdk.StorageClient{
    Bbolt: &sdk.BboltStorage{
        SotoragePath: "./mybolt.db",
        BucketName:   "test",
    },
}
```
#### Redis
```go
storage := &StorageClient{
	Redis: &RedisStorage{
		Address:  "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	},
}
```

### SDK
Create a new SDK instance
```go
avocado := sdk.Avocado{}
err := avocado.New(storage)
if err != nil {
    return err
}
```

### Methods
#### Encrypt and save value
Value can be encrypted and stored using the following methods
```go
encryptedValue, err := avocado.EecryptAndStoreValue([]byte("key1"), []byte("my value"), "/path/to/my_public_key.pem")
```

#### Find and decrypt value
```go
decryptedValue, err := avocado.FindAndDecryptValueBy([]byte("key1"), "/path/to/my_private_key.pem")
```

#### Get all keys
```go
keys, err := avocado.GetAllKeys()
```