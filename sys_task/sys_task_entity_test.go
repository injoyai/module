package sys_task

import (
	"github.com/injoyai/module/sys_corn"
	"github.com/injoyai/module/sys_go"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New(&Config{
		Config: &sys_go.Config{
			DoneSize: 100,
			GoLimit:  100,
			WaitCap:  1000,
		},
		ScriptNum: 10,
	})
	c := &SysTaskCreateReq{
		Group:  "系统",
		Name:   "测试",
		Memo:   "",
		Spec:   sys_corn.NewIntervalSpec(time.Second),
		Script: "print(1)",
	}
	_ = x
	if err := x.PostTask(c); err != nil {
		t.Error(err)
	}
	go func() {
		<-time.After(time.Second * 3)
		c.Script = "print(2)"
	}()

	select {}
}
