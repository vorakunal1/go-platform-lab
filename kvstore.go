package main

import {
	"sync"
	"fmt"
}

type KVstore struct {
	mu sync.RWMutex
	data map[string]interface{}
}

func NewKVStore() *KVstore {
	return &KVstore{data : make(map[string]interface{})}
}

func (s *KVStore) Put(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *KVStore) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RLock()
	val, exists := s.data[key]
	return val, exists
}


func (s *KVStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

