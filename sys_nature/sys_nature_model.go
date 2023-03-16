package sys_nature

import (
	"github.com/injoyai/base/maps"
	"time"
	"xorm.io/xorm"
)

var db *xorm.Engine

type SysNature struct {
	m *maps.Safe
}

func (this *SysNature) Set(group string, key string, val Nature, expiration ...time.Duration) {
	v, _ := this.m.GetOrSetByHandler(group, func() (interface{}, error) { return maps.NewSafe(), nil })
	m := v.(*maps.Safe)
	m.Set(key, val, expiration...)
}

func (this *SysNature) Get(group string, key string) (Nature, bool) {
	v, _ := this.m.GetOrSetByHandler(group, func() (interface{}, error) { return maps.NewSafe(), nil })
	data, has := v.(*maps.Safe).Get(key)
	if !has {
		return nil, false
	}
	return data.(Nature), true
}

func (this *SysNature) MustGet(group, key string) Nature {
	data, _ := this.Get(group, key)
	return data
}

type Nature interface {
	GetName() string       //数据名称
	GetKey() string        //数据唯一标识
	GetValue() interface{} //数据值
	GetType() string       //消息数据类型,string,float,int,bool

	GetLastTime() int64            //最后数据时间,毫秒
	Writable() bool                //是否可写
	WriteValue(value string) error //写入数据
}
