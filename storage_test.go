package sdk

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestBboltStorageGetSetAndDelete(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// create a new key
	key := []byte("my_key")
	value := []byte("my_value")
	storeClient := &StorageClient{
		Bbolt: &BboltStorage{
			SotoragePath: dir + "/mytest.db",
			BucketName:   "test",
		},
	}

	err = storeClient.New()
	if err != nil {
		t.Fatal(err)
	}

	err = storeClient.set(key, value)
	if err != nil {
		t.Fatal(err)
	}

	returned, err := storeClient.get(key)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(returned, value) {
		t.Fatal("Returned should be", string(value))
	}

	// delete the key
	err = storeClient.delete(key)
	if err != nil {
		t.Fatal(err)
	}

	deletedValue, err := storeClient.get(key)
	if err != nil {
		t.Fatal(err)
	}

	if deletedValue != nil {
		t.Fatal("Returned value should be nil")
	}

}

func TestBboltStorageGetAll(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	key := []byte("my_key")
	value := []byte("my_value")
	storeClient := &StorageClient{
		Bbolt: &BboltStorage{
			SotoragePath: dir + "/mytest.db",
			BucketName:   "test",
		},
	}

	err = storeClient.New()
	if err != nil {
		t.Fatal(err)
	}

	err = storeClient.set(key, value)
	if err != nil {
		t.Fatal(err)
	}

	keys, err := storeClient.getAll()
	if err != nil {
		t.Fatal(err)
	}

	isPresent := false
	for i := 0; i < len(keys); i++ {
		if bytes.Equal(keys[i], key) {
			isPresent = true
		}
	}

	if !isPresent {
		t.Fatal("Returned keys should contains ", string(key))
	}
}
