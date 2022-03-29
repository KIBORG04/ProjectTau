package keymap

import (
	"golang.org/x/exp/constraints"
)

type KeyMap[K comparable, V constraints.Ordered] struct {
	Key   K
	Value V
}

func NewKeyMap[K comparable, V constraints.Ordered](key K, value V) *KeyMap[K, V] {
	return &KeyMap[K, V]{
		Key:   key,
		Value: value,
	}
}

func hasKey[K comparable, V constraints.Ordered](keyMap []*KeyMap[K, V], key K) (*KeyMap[K, V], bool) {
	/*
		for _, myMap := range keyMap {
			if myMap.Key == key {
				return &myMap, true
			}
		}
		return nil, false
	*/
	for i := 0; i < len(keyMap); i++ {
		if keyMap[i].Key == key {
			return keyMap[i], true
		}
	}
	return nil, false
}

func AddElem[K comparable, V constraints.Ordered](keyMap []*KeyMap[K, V], key K, value V) []*KeyMap[K, V] {
	findKeyMap, ok := hasKey[K, V](keyMap, key)
	if !ok {
		keyMap = append(keyMap, NewKeyMap[K, V](key, value))
	} else {
		findKeyMap.Value += value
	}
	return keyMap
}

type MyMap[K comparable, V constraints.Ordered] []*KeyMap[K, V]

func (m MyMap[K, V]) Len() int {
	return len(m)
}

func (m MyMap[K, V]) Less(i, j int) bool {
	return m[i].Value < m[j].Value
}

func (m MyMap[K, V]) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
