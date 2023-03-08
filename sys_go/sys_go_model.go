package sys_go

import (
	"context"
	"fmt"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/g"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/cache"
	"github.com/injoyai/conv"
	"time"
)

type Config struct {
	DoneSize int //历史执行保存数量
	GoLimit  int //最大协程数量
}

func newManage(cfg *Config) *Manage {
	m := &Manage{
		queue:   chans.NewQueueFunc(cfg.GoLimit, cfg.GoLimit),
		wait:    make([]*Create, 0),
		running: maps.NewSafe(),
		done:    cache.NewCycle(cfg.DoneSize),
	}
	return m
}

// Manage 协程管理,可视化
type Manage struct {
	queue   *chans.QueueFunc //协程管理
	wait    []*Create        //可选协程函数
	running *maps.Safe       //正在执行协程
	done    *cache.Cycle     //历史协程执行记录
}

// WaitList 等待执行列表
func (this *Manage) WaitList() (list []*Create) {
	return this.wait
}

// RunList 正在执行的列表
func (this *Manage) RunList() (list []*Info) {
	this.running.Range(func(key, value interface{}) bool {
		list = append(list, value.(*Info))
		return true
	})
	return
}

// DoneList 已经执行完的列表
func (this *Manage) DoneList(limit ...int) (list []*Info) {
	for _, v := range this.done.List(limit...) {
		list = append(list, v.(*Info))
	}
	return
}

// Wait 加入到等待执行
func (this *Manage) Wait(c *Create) {
	this.wait = append(this.wait, c)
}

// Run 执行协程
func (this *Manage) Run(c *Create) string {
	return this.run(c.New(), c.Handler)
}

// RunWait 执行协程
func (this *Manage) RunWait(idx int, param ...g.Map) (string, bool) {
	if len(this.wait) > idx {
		c := this.wait[idx]
		return this.run(c.New(param...), c.Handler), true
	}
	return "", false
}

// run 公共执行协程
func (this *Manage) run(info *Info, handler func(ctx context.Context, a *Manage, m Go) error) string {
	this.running.Set(info.Key, info)
	this.queue.Do(func(no int, num int) {
		defer func() {
			this.done.Add(info)
			this.running.Del(info.Key)
		}()
		info.Run(func(ctx context.Context, m Go) error {
			return handler(ctx, this, m)
		})
	})
	return info.Key
}

// Update 更新协程信息
func (this *Manage) Update(u *Update) {
	if v := this.running.MustGet(u.Key); v != nil {
		v.(*Info).Update(u)
	}
}

// Close 关闭正在执行的协程,通过上下文
func (this *Manage) Close(key string) error {
	if v := this.running.MustGet(key); v != nil {
		return v.(*Info).Close()
	}
	return nil
}

// Update 更新信息,界面修改
type Update struct {
	Key   string `json:"key"`   //唯一标识
	Memo  string `json:"memo"`  //备注
	Param g.Map  `json:"param"` //参数
}

// Create 协程新建配置信息
type Create struct {
	Name    string                                           `json:"name"`  //名称
	Memo    string                                           `json:"memo"`  //备注
	Param   g.Map                                            `json:"param"` //参数
	Handler func(ctx context.Context, a *Manage, m Go) error `json:"-"`     //执行函数
}

func (this *Create) New(param ...g.Map) *Info {
	i := &Info{
		Key:    g.UUID(),
		Name:   this.Name,
		Memo:   this.Memo,
		InDate: time.Now(),
		Param:  this.Param,
	}
	if len(param) > 0 {
		i.Param = param[0]
	}
	i.Extend = conv.NewExtend(i)
	return i
}

// Info 协程信息,储存在内存中
type Info struct {
	//基本信息
	Key    string    `json:"key"`    //唯一标识
	Name   string    `json:"name"`   //名称
	Memo   string    `json:"memo"`   //备注
	InDate time.Time `json:"inDate"` //创建时间
	Param  g.Map     `json:"param"`  //参数

	//结果
	Log    []string `json:"log"`    //日志
	Spend  int      `json:"spend"`  //耗时ms
	Succ   bool     `json:"succ"`   //执行是否成功
	Result string   `json:"result"` //执行结果,错误信息

	cancel context.CancelFunc //上下文
	conv.Extend
}

func (this *Info) String() string {
	return fmt.Sprintf("Key:%s Name:%s Param:%s Start:%s",
		this.Key, this.Name, this.Param.Json(), this.InDate.Format(g.TimeLayout))
}

// Update 更新信息
func (this *Info) Update(u *Update) *Info {
	this.Memo = u.Memo
	this.Param = u.Param
	return this
}

// Done 执行结束,赋值
func (this *Info) Done(err error) {
	this.Succ = err == nil
	this.Result = conv.New(err).String("成功")
	this.Spend = int(time.Now().Sub(this.InDate) / 1e6)
	this.cancel = nil
	this.Extend = nil
}

// GetVar 实现接口
func (this *Info) GetVar(key string) *conv.Var {
	return conv.New(this.Param[key])
}

// Print 实现接口,打印日志
func (this *Info) Print(v ...interface{}) {
	msg := fmt.Sprint(v...)
	fmt.Println(msg)
	this.Log = append(this.Log, msg)
}

// Run 协程执行
func (this *Info) Run(fn func(ctx context.Context, g Go) error) {
	ctx, cancel := g.WithCancel()
	this.cancel = cancel
	this.Done(fn(ctx, this))
}

// Close 关闭协程
func (this *Info) Close() error {
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}

type Go interface {
	conv.Extend
	Print(v ...interface{})
}
