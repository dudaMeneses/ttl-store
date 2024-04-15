package main

import (
	"fmt"
	"time"
)

func main() {
	kvs := NewKeyValueStore(2)
	kvs.Set("key1", "value1")

	fmt.Printf("valid values for '%s': %d \n", "key1", kvs.GetValid("key1"))

	time.Sleep(1 * time.Second)

	fmt.Printf("valid values for '%s': %d \n", "key1", kvs.GetValid("key1"))

	time.Sleep(4 * time.Second)

	fmt.Printf("valid values for '%s': %d \n", "key1", kvs.GetValid("key1"))

	kvs.Set("key1", "value1")
	kvs.Set("key1", "value11")
	kvs.Set("key2", "value2")
	kvs.Set("key3", "value3")
	kvs.Set("key3", "value31")
	kvs.Set("key3", "value32")

	fmt.Printf("all valid values: %d \n", kvs.GetAllValid())
}

type Value struct {
	value      string
	expiration time.Time
}

type KeyValueStore struct {
	kvmap map[string][]Value
	ttl   time.Duration
}

func NewKeyValueStore(ttlSec time.Duration) *KeyValueStore {
	return &KeyValueStore{
		kvmap: make(map[string][]Value),
		ttl:   time.Second * ttlSec,
	}
}

func (kvs *KeyValueStore) Set(key string, value string) {
	if list, ok := kvs.kvmap[key]; ok {
		kvs.kvmap[key] = append(list, Value{value: value, expiration: time.Now().Add(kvs.ttl)})
	} else {
		kvs.kvmap[key] = []Value{{value: value, expiration: time.Now().Add(kvs.ttl)}}
	}

	go func() {
		time.Sleep(kvs.ttl)
		kvs.DeleteInvalid(key)
	}()
}

func (kvs *KeyValueStore) GetValid(key string) int {
	count := 0

	if values, ok := kvs.kvmap[key]; ok {
		for _, value := range values {
			if value.expiration.After(time.Now()) {
				count++
			}
		}
	}

	return count
}

func (kvs *KeyValueStore) GetAllValid() int {
	count := 0

	for key := range kvs.kvmap {
		count += kvs.GetValid(key)
	}

	return count
}

func (kvs *KeyValueStore) DeleteInvalid(key string) {
	if values, ok := kvs.kvmap[key]; ok {
		for i, value := range values {
			if value.expiration.Before(time.Now()) {
				kvs.kvmap[key] = append(values[:i], values[i+1:]...)
			}
		}
	}
}
