package v2

import (
	"context"
	"sync"
)

type k int

var contextKey = k(0)

type valuesMap struct {
	store     map[string][]interface{}
	parentCtx context.Context
	lock      sync.RWMutex
}

func newValuesMap(parentCtx context.Context) *valuesMap {
	return &valuesMap{
		store:     make(map[string][]interface{}),
		parentCtx: parentCtx,
	}
}

// Store returns a copy of parent in which the value associated with key is value
func Store[T any](ctx context.Context, key string, values ...T) context.Context {
	vm, ok := ctx.Value(contextKey).(*valuesMap)
	if !ok {
		vm = newValuesMap(ctx)
		ctx = context.WithValue(ctx, contextKey, vm)
	}

	vm.lock.Lock()
	defer vm.lock.Unlock()
	for _, value := range values {
		vm.store[key] = append(vm.store[key], value)
	}

	return ctx
}

// StoreSingleValue returns a copy of parent in which the value associated with key is value
func StoreSingleValue[T any](ctx context.Context, key string, value T) context.Context {
	vm, ok := ctx.Value(contextKey).(*valuesMap)
	if !ok {
		vm = newValuesMap(ctx)
		ctx = context.WithValue(ctx, contextKey, vm)
	}

	vm.lock.Lock()
	defer vm.lock.Unlock()
	vm.store[key] = []interface{}{value}

	return ctx
}

// Load returns the first value stored in the context for the given key.
func Load[T any](ctx context.Context, key string) (value T) {
	values := LoadAll[T](ctx, key)
	if len(values) > 0 {
		return values[0]
	}
	return
}

// LoadAll returns all value associated with this context for key
func LoadAll[T any](ctx context.Context, key string) []T {
	vm, ok := ctx.Value(contextKey).(*valuesMap)
	if !ok {
		return nil
	}

	vm.lock.RLock()
	defer vm.lock.RUnlock()
	values, ok := vm.store[key]
	if !ok {
		return nil
	}

	var result []T
	for _, v := range values {
		value, ok := v.(T)
		if ok {
			result = append(result, value)
		}
	}
	return result
}
