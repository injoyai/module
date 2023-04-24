package sys_nature

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"time"
	"xorm.io/xorm"
)

var db *xorm.Engine

func NewSysNature() *SysNature {
	return &SysNature{
		m: maps.NewSafe(),
	}
}

type SysNature struct {
	m *maps.Safe
}

func (this *SysNature) Set(group string, key string, val INature, expiration ...time.Duration) {
	v, _ := this.m.GetOrSetByHandler(group, func() (interface{}, error) { return maps.NewSafe(), nil })
	m := v.(*maps.Safe)
	m.Set(key, val, expiration...)
}

func (this *SysNature) Get(group string, key string) (INature, bool) {
	v, _ := this.m.GetOrSetByHandler(group, func() (interface{}, error) { return maps.NewSafe(), nil })
	data, has := v.(*maps.Safe).Get(key)
	if !has {
		return nil, false
	}
	return data.(INature), true
}

func (this *SysNature) MustGet(group, key string) INature {
	data, _ := this.Get(group, key)
	return data
}

type INature interface {
	GetName() string       //数据名称
	GetKey() string        //数据唯一标识
	GetValue() interface{} //数据值
	GetType() string       //消息数据类型,string,float,int,bool

	GetLastTime() int64            //最后数据时间,毫秒
	Writable() bool                //是否可写
	WriteValue(value string) error //写入数据
}

type Nature struct {
	Group    string `json:"group"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	LastTime int64  `json:"lastTime"`
}

func (this *Nature) GetGroup() string { return this.Group }

func (this *Nature) GetName() string { return this.Name }

func (this *Nature) GetType() string { return this.Type }

func (this *Nature) GetKey() string { return this.Key }

func (this *Nature) GetValue() interface{} { return this.Value }

func (this *Nature) GetValueString() string { return this.Value }

func (this *Nature) GetLastTime() int64 { return this.LastTime }

func NewNatureNull(group, key string) *Nature {
	return &Nature{Group: group, Key: key}
}

func NewNature(group string, i INature) *Nature {
	return &Nature{
		Group:    group,
		Name:     i.GetName(),
		Type:     i.GetType(),
		Key:      i.GetKey(),
		Value:    conv.String(i.GetValue()),
		LastTime: i.GetLastTime(),
	}
}
