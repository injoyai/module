package sys_task

import (
	"context"
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
		goroute: sys_go.NewManage(cfg.Config),
		script:  sys_script.New(cfg.ScriptNum),
		cron:    sys_corn.New(),
	}
	e.cron.Start()
	return e
}

type Entity struct {
	goroute *sys_go.Manage   //协程管理
	script  *sys_script.Pool //脚本实例
	cron    *sys_corn.Cron   //定时器
}

func (this *Entity) Run(c *Create) error {
	return this.cron.SetTask(c.Name, c.Spec, func() {
		this.goroute.Go(&sys_go.Create{
			Name:  c.Name,
			Memo:  c.Memo,
			Param: nil,
			Handler: func(ctx context.Context, a *sys_go.Manage, m sys_go.Go) error {
				_, err := this.script.Exec(c.Script)
				return err
			},
		})
	})
}

type Create struct {
	Name   string `json:"name"`   //名称
	Memo   string `json:"memo"`   //备注
	Spec   string `json:"spec"`   //表达式
	Script string `json:"script"` //脚本
}

func (this *Create) Handler(goroute *sys_go.Manage, script *sys_script.Pool) {
	goroute.Go(&sys_go.Create{
		Name:  this.Name,
		Memo:  this.Memo,
		Param: nil,
		Handler: func(ctx context.Context, a *sys_go.Manage, m sys_go.Go) error {
			_, err := script.Exec(this.Script)
			return err
		},
	})
}

type Update struct {
	Key string `json:"key"`
	*Create
}

func (this *Entity) Update(u *Update) {
	e := this.cron.GetTask(u.Key)
	_ = e
	//e.
}
