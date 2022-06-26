package main

type Storage struct {
	data map[string]string
}

func (storage *Storage) Get(key string) (value string) {
	value = storage.data[key]
	return
}

func (storage *Storage) Set(key string, value string) {
	storage.data[key] = value
}

func (storage *Storage) Del(key string) {
	delete(storage.data, key)
}
