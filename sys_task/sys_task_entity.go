package sys_task

import (
	"github.com/injoyai/module/sys_corn"
	"github.com/injoyai/module/sys_go"
	"github.com/injoyai/module/sys_script"
)

func New() {

}

type Client struct {
	goroute sys_go.Manage    //协程管理
	script  *sys_script.Pool //脚本实例
	cron    *sys_corn.Cron   //定时器
}
