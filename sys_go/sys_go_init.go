package sys_go

import (
	"errors"
)

var Default = DefaultManage()

func New(cfg *Config) *Entity {
	return newEntity(cfg)
}

func DefaultManage() *Entity {
	m := newEntity(&Config{
		DoneSize: 1000,
		GoLimit:  1000,
		WaitCap:  1000,
	})
	return m
}

var (
	ErrContext = errors.New("上下文关闭")
	ErrTimeout = errors.New("超时")
)
