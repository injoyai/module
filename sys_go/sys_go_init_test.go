package sys_go

import (
	"context"
	"testing"
	"time"
)

func TestDefaultManage(t *testing.T) {
	m := DefaultManage()
	m.Go(Timer)
	m.Go(PrintInfo)
	m.Go(&Create{
		Name:  "测试",
		Param: nil,
		Handler: func(ctx context.Context, a *Manage, m Go) error {
			time.Sleep(time.Second * 10)
			return nil
		},
	})
	select {}
}

func TestManageLimit(t *testing.T) {
	m := NewManage(&Config{
		GoLimit: 2,
	})
	for i := 0; i < 10; i++ {
		m.Go(PrintInfo)
	}
	select {}
}
