package contexthelper_test

import (
	"context"
	"os"
	"testing"

	"gl.fotechwealth.com.local/backend/trade-lib.git/contexthelper"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
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

	ctx = contexthelper.Store(ctx, "requestId", requestId)
	ctx = contexthelper.Store(ctx, "requestContext", requestContext)
	ctx = contexthelper.Store(ctx, "uin", uin)
	ctx = contexthelper.Store(ctx, "pid", pid)

	// 使用标准 context 库完全不受影响
	ctx = context.WithValue(ctx, "son", "2")

	s.Equal(requestId, contexthelper.Load(ctx, "requestId"))
	s.Equal(requestContext, contexthelper.Load(ctx, "requestContext"))
	s.Equal(uin, contexthelper.Load(ctx, "uin"))
	s.Equal(pid, contexthelper.Load(ctx, "pid"))
	s.Equal("1", ctx.Value("parent"))
	s.Equal("2", ctx.Value("son"))

	// 存入多个值映射到相同 key
	ctx = contexthelper.Store(ctx, "uin", "2")
	ctx = contexthelper.Store(ctx, "uin", "3")

	s.Equal([]interface{}{"1", "2", "3"}, contexthelper.LoadAll(ctx, "uin"))
}

func TestAll(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func TestStoreSingleValue(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ctx := context.Background()
		ctx = contexthelper.Store(ctx, "requestId", uuid.NewV4().String())
		t.Logf("ctx value:%+v", contexthelper.LoadAll(ctx, "requestId"))

		ctx = contexthelper.Store(ctx, "requestId", uuid.NewV4().String())
		t.Logf("ctx value:%+v", contexthelper.LoadAll(ctx, "requestId"))

		ctx = contexthelper.Store(ctx, "requestId", uuid.NewV4().String())
		t.Logf("ctx value:%+v", contexthelper.LoadAll(ctx, "requestId"))
	})

	t.Run("", func(t *testing.T) {
		ctx := context.Background()
		ctx = contexthelper.StoreSingleValue(ctx, "requestId", uuid.NewV4().String())
		t.Logf("ctx value:%+v", contexthelper.LoadAll(ctx, "requestId"))

		ctx = contexthelper.StoreSingleValue(ctx, "requestId", uuid.NewV4().String())
		t.Logf("ctx value:%+v", contexthelper.LoadAll(ctx, "requestId"))

		ctx = contexthelper.StoreSingleValue(ctx, "requestId", uuid.NewV4().String())
		t.Logf("ctx value:%+v", contexthelper.LoadAll(ctx, "requestId"))
	})
}