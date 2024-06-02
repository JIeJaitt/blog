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
