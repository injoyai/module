package sys_task

import (
	"errors"
	"github.com/injoyai/module/sys_go"
	"github.com/injoyai/module/sys_script"
	"github.com/robfig/cron/v3"
	"xorm.io/xorm"
)

type Config struct {
	*sys_go.Config
	ScriptNum int
	DataDir   string
}

func New(cfg *Config) *Entity {
	e := &Entity{
		goroute: sys_go.New(cfg.Config),
		script:  sys_script.New(cfg.ScriptNum),
		cron:    cron.New(cron.WithSeconds()),
	}
	e.cron.Start()
	return e
}

type Entity struct {
	goroute *sys_go.Entity   //协程管理
	script  *sys_script.Pool //脚本实例
	cron    *cron.Cron       //定时器
	db      *xorm.Engine     //数据库
}

//func (this *Entity) Loading() error {
//
//}

func (this *Entity) GetTaskAll() (list []*SysTask, err error) {
	err = this.db.Find(&list)
	return
}

func (this *Entity) GetTaskList() {

}

func (this *Entity) GetTask(id int64) (*SysTask, error) {
	data := new(SysTask)
	has, err := this.db.Where("ID=?", id).Get(data)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("任务不存在")
	}
	return data, nil
}

func (this *Entity) PostTask(req *SysTaskCreateReq) error {
	data, _, err := req.New()
	if err != nil {
		return err
	}
	_, err = this.db.Insert(data)
	if err == nil {

	}
	return err
}

func (this *Entity) PutTask(req *SysTaskUpdateReq) error {
	data, cols, err := req.New()
	if err != nil {
		return err
	}
	_, err = this.db.Where("ID=?", req.ID).Cols(cols).Update(data)
	if err == nil {

	}
	return err
}

func (this *Entity) DelTask(id int64) error {
	_, err := this.db.Where("ID=?", id).Delete(new(SysTask))
	if err == nil {
		//this.cron.Remove()
	}
	return err
}

//func (this *Entity) run() error {
//	return this.cron.AddFunc().SetTask(key, c.Spec, func() {
//		this.goroute.Go(&sys_go.Create{
//			Group: c.Group,
//			Name:  c.Name,
//			Memo:  c.Memo,
//			Handler: func(ctx context.Context, a *sys_go.Entity, m sys_go.Param) error {
//				_, err := this.script.Exec(c.Script)
//				return err
//			},
//		})
//	}, c)
//}

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

//func (this *Entity) Update(u *Update) {
//	e := this.cron.GetTask(u.Key)
//	if e != nil {
//		e.Data.(*Create).Update(u)
//	}
//}
