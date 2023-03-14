package sys_go

import (
	"errors"
)

func NewManage(cfg *Config) *Manage {
	return newManage(cfg)
}

func DefaultManage() *Manage {
	m := newManage(&Config{
		DoneSize: 100,
		GoLimit:  100,
		WaitCap:  100,
	})
	return m
}

var (
	ErrContext = errors.New("上下文关闭")
	ErrTimeout = errors.New("超时")
)
