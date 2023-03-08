package sys_go

import (
	"github.com/injoyai/base/g"
	"testing"
)

func TestDefaultManage(t *testing.T) {
	m := DefaultManage()
	m.RunWait(0, g.Map{"text": "test msg"})
	m.RunWait(1)
	select {}
}
