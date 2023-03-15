package sys_task

import (
	"context"
	"errors"
	"github.com/injoyai/base/g"
	"github.com/injoyai/module/sys_corn"
	"github.com/injoyai/module/sys_go"
	"github.com/injoyai/module/sys_script"
)

type Config struct {
	*sys_go.Config
	ScriptNum int
}

func New(cfg *Config) *Entity {
	e := &Entity{
		goroute: sys_go.New(cfg.Config),
		script:  sys_script.New(cfg.ScriptNum),
		cron:    sys_corn.New(),
	}
	e.cron.Start()
	return e
}

type Entity struct {
	goroute *sys_go.Entity   //协程管理
	script  *sys_script.Pool //脚本实例
	cron    *sys_corn.Cron   //定时器
}

func (this *Entity) GetTaskAll() (list []*Info) {
	for _, v := range this.cron.GetTaskAll() {
		list = append(list, &Info{
			Key:    v.Key,
			Create: v.Data.(*Create),
		})
	}
	return
}

func (this *Entity) GetTask(key string) (*Info, error) {
	e := this.cron.GetTask(key)
	if e == nil {
		return nil, errors.New("任务不存在")
	}
	return &Info{
		Key:    e.Key,
		Create: e.Data.(*Create),
	}, nil
}

func (this *Entity) PostTask(c *Create) error {
	key := g.UUID()
	return this.cron.SetTask(key, c.Spec, func() {
		this.goroute.Go(&sys_go.Create{
			Group: c.Group,
			Name:  c.Name,
			Memo:  c.Memo,
			Handler: func(ctx context.Context, a *sys_go.Entity, m sys_go.Param) error {
				_, err := this.script.Exec(c.Script)
				return err
			},
		})
	}, c)
}

func (this *Entity) PutTask(u *Update) error {
	e, err := this.GetTask(u.Key)
	if err != nil {
		return err
	}
	e.Create.Update(u)
	return nil
}

func (this *Entity) DelTask(key string) error {
	this.cron.DelTask(key)
	return nil
}

type Info struct {
	Key string `json:"key"`
	*Create
}

type Create struct {
	Group  string `json:"group"`  //分组
	Name   string `json:"name"`   //名称
	Memo   string `json:"memo"`   //备注
	Spec   string `json:"spec"`   //表达式
	Script string `json:"script"` //脚本
}

func (this *Create) Update(u *Update) {
	this.Group = u.Group
	this.Name = u.Name
	this.Memo = u.Memo
	this.Script = u.Script
}

type Update struct {
	Key    string `json:"key"`    //唯一标识
	Group  string `json:"group"`  //分组
	Name   string `json:"name"`   //名称
	Memo   string `json:"memo"`   //备注
	Spec   string `json:"spec"`   //表达式
	Script string `json:"script"` //脚本
}

func (this *Entity) Update(u *Update) {
	e := this.cron.GetTask(u.Key)
	if e != nil {
		e.Data.(*Create).Update(u)
	}
}
