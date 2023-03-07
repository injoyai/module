package sys_go

import (
	"context"
	"fmt"
	"github.com/injoyai/conv"
	"time"
)

// GoManage 协程管理,可视化
type GoManage struct {
	Key string `json:"key"`
}

func NewGoInfo() *GoInfo {
	return &GoInfo{}
}

// GoInfo 协程信息,储存在内存中
type GoInfo struct {
	Name   string                 `json:"name"`   //名称
	Memo   string                 `json:"memo"`   //备注
	Param  map[string]interface{} `json:"param"`  //参数
	Log    []string               `json:"log"`    //日志
	Succ   bool                   `json:"succ"`   //执行是否成功
	Result string                 `json:"result"` //执行结果
	InDate time.Time              `json:"inDate"` //创建时间
	cancel context.CancelFunc     //上下文
	conv.Extend
}

func (this *GoInfo) GetVar(key string) *conv.Var {
	return conv.New(this.Param[key])
}

func (this *GoInfo) Debug(v ...interface{}) {
	this.Log = append(this.Log, fmt.Sprint(v...))
}

type IGo interface {
	conv.Extend
	Debug(v ...interface{})
}

// Run 协程执行
func (this *GoInfo) Run(fn func(ctx context.Context, m IGo) error) {
	ctx, cancel := context.WithCancel(context.Background())
	this.cancel = cancel
	go func(ctx context.Context) {
		err := fn(ctx, this)
		this.Succ = err == nil
		this.Result = conv.New(err).String("成功")
	}(ctx)
}

// Close 关闭协程
func (this *GoInfo) Close() error {
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}
