package store

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type ConcurrentMap[K comparable, V any] struct {
	sync.RWMutex
	Map map[K]V
}

type Store[K comparable, V any] struct {
	logger TransactionLogger
	ConcurrentMap[K, V]
}

var ErrorNoSuchKey = errors.New("no such key")

var store Store[string, string]

func Put(key, value string) error {
	store.Lock()
	defer store.Unlock()

	store.Map[key] = value
	store.logger.WritePut(key, value)

	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	val, ok := store.Map[key]
	if !ok {
		return "", ErrorNoSuchKey
	}
	return val, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.Map, key)
	store.logger.WriteDelete(key)

	return nil
}

func InitializeStore(filename string) error {
	log.Printf("Initializing store from transaction log: %s", filename)

	logger, err := NewFileTranscationLogger(filename)
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}

	// init store
	store = Store[string, string]{
		logger,
		ConcurrentMap[string, string]{Map: make(map[string]string)},
	}

	// Read the transaction log and restore state
	events, errors := logger.ReadEvents()

	e := Event{}
	ok := true

	for ok && err == nil {
		select {
		case err, ok = <-errors: // retrieve errors
		case e, ok = <-events: // Retrieve events
			switch e.EventType {
			case EventDelete:
				log.Printf("Sequence: %d Event: DELETE Key: %s", e.Sequence, e.Key)
				store.Lock()
				delete(store.Map, e.Key)
				store.Unlock()
				err = nil

			case EventPut:
				log.Printf("Sequence: %d Event: PUT Key: %s Value: %s", e.Sequence, e.Key, e.Value)
				store.Lock()
				store.Map[e.Key] = e.Value
				store.Unlock()
				err = nil
			}
		}
	}

	if err == nil {
		// initialize transcation log for writing new events
		log.Printf("Start transaction logger")
		logger.Run()
	}
	return err
}
