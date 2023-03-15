package sys_go

import (
	"context"
	"fmt"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/g"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/cache"
	"github.com/injoyai/conv"
	"sync/atomic"
	"time"
)

type Config struct {
	DoneSize int  //历史执行保存数量
	GoLimit  uint //最大协程数量
	WaitCap  int  //等待队列长度
}

func newEntity(cfg *Config) *Entity {
	if cfg.GoLimit == 0 {
		cfg.GoLimit = 1000
	}
	if cfg.WaitCap <= 0 {
		cfg.WaitCap = 1000
	}
	m := &Entity{
		limit:    chans.NewWaitLimit(cfg.GoLimit),
		waitList: make([]*WaitInfo, 0),
		waitChan: make(chan *WaitInfo, cfg.WaitCap),
		running:  maps.NewSafe(),
		done:     cache.NewCycle(cfg.DoneSize),
	}
	go m.run()
	return m
}

// Entity 协程管理,可视化
type Entity struct {
	cfg        *Config          //配置信息
	limit      *chans.WaitLimit //协程管理
	waitList   []*WaitInfo      //等待执行的协程
	waitChan   chan *WaitInfo   //等待执行的通道
	running    *maps.Safe       //正在执行协程
	runningNum int32            //正在执行的数量
	done       *cache.Cycle     //历史协程执行记录
}

// WaitList 等待列表
func (this *Entity) WaitList() []*WaitInfo {
	return this.waitList
}

// RunNum 释放协程的数量
func (this *Entity) RunNum() int {
	return int(atomic.LoadInt32(&this.runningNum))
}

// RunList 正在执行的列表
func (this *Entity) RunList() (list []*Info) {
	this.running.Range(func(key, value interface{}) bool {
		list = append(list, value.(*Info))
		return true
	})
	return
}

// DoneList 已经执行完的列表
func (this *Entity) DoneList(limit ...int) (list []*Info) {
	for _, v := range this.done.List(limit...) {
		list = append(list, v.(*Info))
	}
	return
}

// Go 执行协程
func (this *Entity) Go(c *Create) *WaitInfo {
	w := c.New()
	this.waitChan <- w
	this.waitList = append(this.waitList, w)
	return w
}

// run 公共执行协程
func (this *Entity) run() {
	for {
		this.limit.Add()
		c := <-this.waitChan
		this.waitList = this.waitList[1:]
		info := c.New()
		this.running.Set(info.Key, info)
		atomic.AddInt32(&this.runningNum, 1)
		go func(info *Info) {
			defer func() {
				atomic.AddInt32(&this.runningNum, -1)
				this.done.Add(info)
				this.running.Del(info.Key)
				this.limit.Done()
			}()
			info.Run(func(ctx context.Context, p Param) (err error) {
				defer g.Recover(&err)
				return c.Handler(ctx, this, p)
			})
		}(info)
	}
}

// Update 更新协程信息
func (this *Entity) Update(u *Update) {
	if v := this.running.MustGet(u.Key); v != nil {
		v.(*Info).Update(u)
	}
}

// Close 关闭正在执行的协程,通过上下文
func (this *Entity) Close(key string) error {
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
	Group   string                                              `json:"group"` //分组
	Name    string                                              `json:"name"`  //名称
	Memo    string                                              `json:"memo"`  //备注
	Param   g.Map                                               `json:"param"` //参数
	Handler func(ctx context.Context, a *Entity, p Param) error `json:"-"`     //执行函数
}

func (this *Create) New() *WaitInfo {
	return &WaitInfo{
		Key:    g.UUID(),
		Create: this,
	}
}

type WaitInfo struct {
	Key string `json:"key"` //唯一标识
	*Create
}

func (this *WaitInfo) New(param ...g.Map) *Info {
	i := &Info{
		Key:     this.Key,
		Group:   this.Group,
		Name:    this.Name,
		Memo:    this.Memo,
		RunDate: time.Now(),
		Param:   this.Param,
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
	Key     string    `json:"key"`    //唯一标识
	Group   string    `json:"group"`  //分组信息
	Name    string    `json:"name"`   //名称
	Memo    string    `json:"memo"`   //备注
	Param   g.Map     `json:"param"`  //参数
	RunDate time.Time `json:"inDate"` //开始时间

	//结果
	Log      []string  `json:"log"`      //日志
	Spend    int       `json:"spend"`    //耗时ms
	Succ     bool      `json:"succ"`     //执行是否成功
	Result   string    `json:"result"`   //执行结果,错误信息
	DoneDate time.Time `json:"doneDate"` //结束时间

	cancel context.CancelFunc //上下文
	conv.Extend
}

func (this *Info) String() string {
	return fmt.Sprintf("Key:%s  Name:%s  Param:%s  Start:%s",
		this.Key, this.Name, this.Param.Json(), this.RunDate.Format(g.TimeLayout))
}

// Update 更新信息
func (this *Info) Update(u *Update) {
	this.Memo = u.Memo
	this.Param = u.Param
}

// Done 执行结束,赋值
func (this *Info) Done(err error) {
	this.Succ = err == nil
	this.Result = conv.New(err).String("成功")
	this.DoneDate = time.Now()
	this.Spend = int(this.DoneDate.Sub(this.RunDate) / 1e6)
	this.cancel = nil
	this.Extend = nil
}

// GetVar 实现接口
func (this *Info) GetVar(key string) *conv.Var {
	return conv.New(this.Param[key])
}

// Print 实现接口,打印日志
func (this *Info) Print(v ...interface{}) {
	this.Log = append(this.Log, fmt.Sprint(v...))
}

// Run 协程执行
func (this *Info) Run(fn func(ctx context.Context, p Param) error) {
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

type Param interface {
	conv.Extend
	Print(v ...interface{})
}
