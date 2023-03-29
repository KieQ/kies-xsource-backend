package utils

import (
	"fmt"
	"math/rand"
	"reflect"
)

func AddToMapIfNotNil[K comparable, V any, M map[K]any](m M, v *V, k K) {
	if v == nil {
		return
	}
	m[k] = *v
}


func GetFromAnyMap[T any, K comparable](m map[K]any, k K) (T, error) {
	var result T
	var ok bool
	if val, exist := m[k]; !exist {
		return result, fmt.Errorf("key %v does not exist", k)
	} else if result, ok = val.(T); !ok {
		return result, fmt.Errorf("value does not has type %v, it's %v", reflect.TypeOf(result), reflect.TypeOf(val))
	} else {
		return result, nil
	}
}

func Sample[T any](samples []T, count int) []T {
	if len(samples) <= count {
		return samples
	}
	result := make([]T, 0, count)
	for _, v := range rand.Perm(len(samples))[:count] {
		result = append(result, samples[v])
	}
	return result
}

func UniqueBy[T any, K comparable](items []T, f func(T) K) []T {
	m := make(map[K]T, len(items))
	for _, item := range items {
		m[f(item)] = item
	}
	result := make([]T, 0, len(m))
	for _, item := range m {
		result = append(result, item)
	}
	return result
}
