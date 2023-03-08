package sys_go

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func NewManage(cfg *Config) *Manage {
	return newManage(cfg)
}

func DefaultManage() *Manage {
	m := newManage(&Config{
		DoneSize: 100,
		GoLimit:  10,
	})
	m.Wait(Timer)
	m.Wait(PrintInfo)
	return m
}

var (
	Timer = &Create{
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
					m.Print(m.GetString("text"))
				}
			}
		},
	}
	PrintInfo = &Create{
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
					fmt.Println("=========================================================")
					for _, v := range a.RunList() {
						fmt.Println(v)
						fmt.Println("Log:")
						for _, k := range v.Log {
							fmt.Println("\t", k)
						}
					}
				}
			}
		},
	}
)

var (
	ErrContext = errors.New("上下文关闭")
	ErrTimeout = errors.New("超时")
)
