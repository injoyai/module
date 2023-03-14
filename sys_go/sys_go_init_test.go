package sys_go

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	testTimer = &Create{
		Name:  "定时器",
		Memo:  "",
		Param: nil,
		Handler: func(ctx context.Context, a *Manage, m Go) error {
			second := m.GetSecond("second", 1)
			timer := time.NewTimer(second)
			defer timer.Stop()
			for {
				timer.Reset(second)
				select {
				case <-ctx.Done():
					return ErrContext
				case <-timer.C:
					m.Print(m.GetString("text", "text"))
				}
			}
		},
	}
	testPrint = &Create{
		Name:  "打印信息",
		Memo:  "",
		Param: nil,
		Handler: func(ctx context.Context, a *Manage, m Go) error {
			second := m.GetSecond("second", 1)
			timer := time.NewTimer(second)
			defer timer.Stop()
			for {
				timer.Reset(second)
				select {
				case <-ctx.Done():
					return ErrContext
				case <-timer.C:
					fmt.Println("===============================Run===========================")
					for _, v := range a.RunList() {
						fmt.Println(v)
					}
					fmt.Println("===============================Done==========================")
					for _, v := range a.DoneList() {
						fmt.Println(v)
					}
				}
			}
		},
	}
)

func TestDefaultManage(t *testing.T) {
	m := DefaultManage()
	m.Go(testTimer)
	m.Go(testPrint)
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
		m.Go(testPrint)
	}
	select {}
}
