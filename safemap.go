package safemap

import (
	"encoding/json"
	"sync"
)

type Map[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok = m.m[key]

	return value, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.m[key] = value

	return
}

func (m *Map[K, V]) All() map[K]V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	all := make(map[K]V, len(m.m))
	for k, v := range m.m {
		all[k] = v
	}

	return all
}

func (m *Map[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.m, key)
}

func (m *Map[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.m = make(map[K]V)
}

func (m *Map[K, V]) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.m)
}

func (m *Map[K, V]) Random() (key K, value V, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for k, v := range m.m {
		return k, v, true
	}

	return key, value, false
}

func (m *Map[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var keys []K
	for k := range m.m {
		keys = append(keys, k)
	}

	return keys
}

func (m *Map[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var values []V
	for _, v := range m.m {
		values = append(values, v)
	}

	return values
}

func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return json.Marshal(m.m)
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	v := make(map[K]V)
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	m.mutex.Lock()
	m.m = v
	m.mutex.Unlock()

	return nil
}

func (m *Map[K, V]) MarshalText() ([]byte, error) {
	return m.MarshalJSON()
}

func (m *Map[K, V]) UnmarshalText(data []byte) error {
	return m.UnmarshalJSON(data)
}
