package contexthelper

import (
	"context"

	constant "gl.fotechwealth.com.local/backend/trade-lib.git/contexthelper/constant"
)

// use share contextKey
var ContextKey = constant.SharedContextKey

// Store returns a copy of parent in which the value associated with key is value
func Store(ctx context.Context, key string, value interface{}) context.Context {
	var m map[string][]interface{}
	i := ctx.Value(ContextKey)
	if i == nil {
		m = map[string][]interface{}{}
		ctx = context.WithValue(ctx, ContextKey, m)
	} else {
		m = i.(map[string][]interface{})
	}

	m[key] = append(m[key], value)
	return ctx
}

// StoreSingleValue returns a copy of parent in which the value associated with key is value
func StoreSingleValue(ctx context.Context, key string, value interface{}) context.Context {
	var m map[string][]interface{}
	i := ctx.Value(ContextKey)
	if i == nil {
		m = map[string][]interface{}{}
		ctx = context.WithValue(ctx, ContextKey, m)
	} else {
		m = i.(map[string][]interface{})
	}

	m[key] = []interface{}{value}
	return ctx
}

// Load returns the first value associated with this context for key
func Load(ctx context.Context, key string) interface{} {
	values := LoadAll(ctx, key)
	if len(values) < 1 {
		return nil
	}

	return values[0]
}

// LoadAll returns all value associated with this context for key
func LoadAll(ctx context.Context, key string) []interface{} {
	i := ctx.Value(ContextKey)
	if i == nil {
		return nil
	}

	return i.(map[string][]interface{})[key]
}
