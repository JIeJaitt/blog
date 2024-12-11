package v2_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gl.fotechwealth.com.local/backend/trade-lib.git/contexthelper"
	v2 "gl.fotechwealth.com.local/backend/trade-lib.git/contexthelper/v2"
)

type TestSuite struct {
	suite.Suite
}

func (s *TestSuite) TestRun() {
	ctx := context.WithValue(context.TODO(), "parent", "1")
	requestId := uuid.NewV4().String()
	requestContext := "POST:/ping"
	uin := "1"
	pid := os.Getpid()

	ctx = v2.Store(ctx, "requestId", requestId)
	ctx = v2.Store(ctx, "requestContext", requestContext)
	ctx = v2.Store(ctx, "uin", uin)
	ctx = v2.Store(ctx, "pid", pid)

	// 使用标准 context 库完全不受影响
	ctx = context.WithValue(ctx, "son", "2")

	s.Equal(requestId, v2.Load[string](ctx, "requestId"))
	s.Equal(requestContext, v2.Load[string](ctx, "requestContext"))
	s.Equal(uin, v2.Load[string](ctx, "uin"))
	s.Equal(pid, v2.Load[int](ctx, "pid"))
	s.Equal("1", ctx.Value("parent"))
	s.Equal("2", ctx.Value("son"))

	// 存入多个值映射到相同 key
	ctx = v2.Store(ctx, "uin", "2")
	ctx = v2.Store(ctx, "uin", "3")

	s.Equal([]string{"1", "2", "3"}, v2.LoadAll[string](ctx, "uin"))
}

func TestAll(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func TestStoreMultipleValues(t *testing.T) {
	ctx := context.Background()
	requestId := uuid.NewV4().String()

	// 存储多个值到相同的键
	ctx = v2.Store(ctx, "requestId", requestId, "anotherId", "yetAnotherId")

	// 从context中加载所有的值
	values := v2.LoadAll[string](ctx, "requestId")
	assert.Equal(t, 3, len(values))
	assert.Equal(t, requestId, values[0])
	assert.Equal(t, "anotherId", values[1])
	assert.Equal(t, "yetAnotherId", values[2])
}

func TestStoreSingleValue(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ctx := context.Background()
		requestId := uuid.NewV4().String()
		ctx = v2.StoreSingleValue(ctx, "requestId", requestId)
		t.Logf("ctx value:%+v", v2.Load[string](ctx, "requestId"))

		newRequestId := uuid.NewV4().String()
		ctx = v2.StoreSingleValue(ctx, "requestId", newRequestId)
		t.Logf("ctx value:%+v", v2.Load[string](ctx, "requestId"))
	})
}

func TestWriteDoesNotAffectParentContext(t *testing.T) {
	parentCtx := context.Background()
	childCtx := v2.Store(parentCtx, "childKey", "childValue")

	// 子context应有新值
	childValue := v2.Load[string](childCtx, "childKey")
	assert.Equal(t, "childValue", childValue)

	// 父context不应有子context的值
	parentValue := v2.Load[string](parentCtx, "childKey")
	assert.Equal(t, "", parentValue)
}

// TestCascadingReadFromParentContext
// TODO 优化
func TestCascadingReadFromParentContext(t *testing.T) {
	parentCtx := context.Background()
	parentCtx = v2.Store(parentCtx, "parentKey", "parentValue")
	childCtx := v2.Store(parentCtx, "childKey", "childValue")

	// 子context应有新值
	childValue := v2.Load[string](childCtx, "childKey")
	assert.Equal(t, "childValue", childValue)

	// 子context应级联读取父context的值
	parentValue := v2.Load[string](childCtx, "parentKey")
	assert.Equal(t, "parentValue", parentValue)
}

// 测试并发功能
func TestConcurrentReadWrite(t *testing.T) {
	ctx := context.Background()

	// 创建一个临时上下文，以便在写入完成后切换回
	tmpCtx := context.Background()

	// 并发写入
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		uniqueKey := fmt.Sprintf("concurrentKey_%d", i)

		wg.Add(1)
		go func(i int, key string) {
			defer wg.Done()
			tmpCtx = v2.Store(tmpCtx, key, i)
		}(i, uniqueKey)
	}

	// 确保所有写入操作完成
	wg.Wait()

	// 将写入完成后的临时上下文赋值给原始上下文
	ctx = tmpCtx

	// 并发读取
	wg = sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		uniqueKey := fmt.Sprintf("concurrentKey_%d", i)

		wg.Add(1)
		go func(i int, key string) {
			defer wg.Done()
			value := v2.Load[int](ctx, key)
			if value != i {
				t.Errorf("键 %s 的值应该是 %d，但得到的是 %v", key, i, value)
			}
		}(i, uniqueKey)
	}

	// 等待所有读取完成
	wg.Wait()

	// 如果没有出现任何 Error 或 Fatal，那么测试就是成功的
	t.Log("TestConcurrentReadWrite passed")
}

// TestConcurrentReadWhileWrite 边读边写，没有报错就是成功
func TestConcurrentReadWhileWrite(t *testing.T) {
	originalCtx := context.Background()
	var ctx context.Context
	var wg sync.WaitGroup
	key := "concurrentKey"

	// 定义写入函数
	writeFunc := func(value string) {
		ctx = v2.Store(originalCtx, key, value)
	}

	// 定义读取函数
	readFunc := func() string {
		return v2.Load[string](ctx, key)
	}

	// 启动多个写入 goroutines
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value := fmt.Sprintf("value%d", i)
			writeFunc(value)
		}(i)
	}

	// 启动多个读取 goroutines
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			readValue := readFunc()
			t.Log("Read value:", readValue)
			// 这里不使用断言，因为读取到的值可能会在写入时改变
		}()
	}

	// 等待所有 goroutines 完成
	wg.Wait()
}

func TestCompatibility(t *testing.T) {
	ctx := context.Background()

	// 使用旧版本存储值
	ctx = contexthelper.Store(ctx, "testKey", "testValue")

	// 尝试使用新版本读取值
	retrievedValue := v2.Load[string](ctx, "testKey")
	assert.Equal(t, "testValue", retrievedValue, "The value retrieved by new version should match the value stored by old version.")
}

// TestWriteDoesNotAffectParent 测试写入只影响当前上下文
func TestWriteDoesNotAffectParent(t *testing.T) {
	parentCtx := context.Background()
	childCtx := context.WithValue(parentCtx, v2.ContextKey, v2.NewValuesMap(parentCtx))

	// 在子上下文中存储值
	v2.Store(childCtx, "childKey", "childValue")

	// 确认父上下文没有被影响
	parentValue := v2.Load[string](parentCtx, "childKey")
	assert.Empty(t, parentValue, "Parent context should not have the child's value")
}

// TestCascadingRead 测试上下文级联读取
func TestCascadingReadFromParentContextV2(t *testing.T) {
	// 创建父上下文并存储值
	parentCtx := context.Background()
	parentCtx = v2.Store(parentCtx, "parentKey", "parentValue")

	// 验证父上下文中是否正确存储了值
	parentValue := v2.Load[string](parentCtx, "parentKey")
	assert.Equal(t, "parentValue", parentValue, "Value should be stored in parent context")

	// 创建子上下文，并传递父上下文引用
	childCtx := context.WithValue(parentCtx, v2.ContextKey, v2.NewValuesMap(parentCtx))

	// 子上下文应该能够级联读取父上下文的值
	childValue := v2.Load[string](childCtx, "parentKey")
	assert.Equal(t, "parentValue", childValue, "Should retrieve value from parent context")

	// 在子上下文中存储值
	childCtx = v2.Store(childCtx, "childKey", "childValue")

	// 确保子上下文可以读取自己存储的值
	valueFromChild := v2.Load[string](childCtx, "childKey")
	assert.Equal(t, "childValue", valueFromChild, "Should retrieve value stored in child context")

	// 确保父上下文不能读取子上下文存储的值
	valueFromParent := v2.Load[string](parentCtx, "childKey")
	assert.Equal(t, "", valueFromParent, "Parent context should not retrieve value stored in child context")
}
