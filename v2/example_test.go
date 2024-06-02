package v2_test

import (
	"context"
	"fmt"
	v2 "gl.fotechwealth.com.local/backend/trade-lib.git/contexthelper/v2"
)

func ExampleStore() {
	ctx := context.Background()
	ctx = v2.Store(ctx, "key1", "value1")
	value := v2.Load[string](ctx, "key1")
	fmt.Println(value)
	// Output: value1
}

func ExampleStore_multipleValues() {
	ctx := context.Background()
	ctx = v2.Store(ctx, "key1", "value1", "value2", "value3")
	values := v2.LoadAll[string](ctx, "key1")
	fmt.Println(values)
	// Output: [value1 value2 value3]
}

func ExampleStoreSingleValue() {
	ctx := context.Background()
	ctx = v2.StoreSingleValue(ctx, "key1", "value1")
	ctx = v2.StoreSingleValue(ctx, "key1", "value2")
	ctx = v2.StoreSingleValue(ctx, "key1", "value3")
	value := v2.Load[string](ctx, "key1")
	fmt.Println(value)
	// Output: value3
}

func ExampleLoad() {
	ctx := context.Background()
	ctx = v2.Store(ctx, "key1", "value1")
	ctx = v2.Store(ctx, "key1", "value2")
	ctx = v2.Store(ctx, "key1", "value3")
	value := v2.Load[string](ctx, "key1")
	fmt.Println(value)
	// Output: value1
}

func ExampleLoadAll() {
	ctx := context.Background()
	ctx = v2.Store(ctx, "key1", "value1")
	ctx = v2.Store(ctx, "key1", "value2")
	ctx = v2.Store(ctx, "key1", "value3")
	values := v2.LoadAll[string](ctx, "key1")
	fmt.Println(values)
	// Output: [value1 value2 value3]
}
