package main

import (
	"context"
	"fmt"

	v2 "gl.fotechwealth.com.local/backend/trade-lib.git/contexthelper/v2"
)

func main() {
	ctx := context.Background()

	// 假设这是旧版本的存储方式
	// 直接创建一个map并存储一些值
	oldMap := make(map[string][]interface{})
	oldMap["testKey"] = append(oldMap["testKey"], "testValue")

	// 使用旧的contextKey手动设置这个map到上下文中
	ctx = context.WithValue(ctx, v2.ContextKey, oldMap)

	// 使用新版本的LoadAll尝试读取这个值
	retrievedValues := v2.LoadAll[string](ctx, "testKey")
	if len(retrievedValues) > 0 {
		fmt.Println("Retrieved values:", retrievedValues)
	} else {
		fmt.Println("No values retrieved.")
	}
}
