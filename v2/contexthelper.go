package v2

import (
	"context"
	"sync"

	constant "gl.fotechwealth.com.local/backend/trade-lib.git/contexthelper/constant"
)

// use share contextKey
var ContextKey = constant.SharedContextKey

type valuesMap struct {
	store     map[string][]interface{}
	parentCtx context.Context
	lock      sync.RWMutex
}

func NewValuesMap(parentCtx context.Context) *valuesMap {
	return &valuesMap{
		store:     make(map[string][]interface{}),
		parentCtx: parentCtx,
	}
}

// Store returns a copy of parent in which the value associated with key is value
func Store[T any](ctx context.Context, key string, values ...T) context.Context {
	vm, ok := ctx.Value(ContextKey).(*valuesMap)
	if !ok {
		vm = NewValuesMap(ctx)
		ctx = context.WithValue(ctx, ContextKey, vm)
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
	vm, ok := ctx.Value(ContextKey).(*valuesMap)
	if !ok {
		vm = NewValuesMap(ctx)
		ctx = context.WithValue(ctx, ContextKey, vm)
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
	vm, ok := ctx.Value(ContextKey).(*valuesMap)
	if !ok {
		// handle old version
		oldMap, ok := ctx.Value(ContextKey).(map[string][]interface{})
		if ok {
			if oldValues, found := oldMap[key]; found {
				var result []T
				for _, v := range oldValues {
					if value, ok := v.(T); ok {
						result = append(result, value)
					}
				}
				return result
			}
		}
		return nil
	}

	vm.lock.RLock()
	defer vm.lock.RUnlock()
	
	values, ok := vm.store[key]
	if !ok || len(values) == 0 {
		// 递归查找父上下文
		if vm.parentCtx != nil {
			return LoadAll[T](vm.parentCtx, key)
		}
		return nil
	}

	var result []T
	for _, v := range values {
		if value, ok := v.(T); ok {
			result = append(result, value)
		}
	}
	return result
}
