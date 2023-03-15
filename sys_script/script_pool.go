package sys_script

import (
	"fmt"
	"github.com/injoyai/conv"
)

func New(num int) *Pool {
	p := &Pool{
		free: make(chan *Client, num),
	}
	for i := 0; i < num; i++ {
		c := NewClient()
		p.all = append(p.all, c)
		p.free <- c
	}
	p.SetFunc("len", funcLen)
	p.SetFunc("base64ToString", funcBase64ToString)
	p.SetFunc("base64ToBytes", funcBase64ToBytes)
	p.SetFunc("hexToBytes", funcHEXToBytes)
	p.SetFunc("hexToString", funcHEXToString)
	p.SetFunc("getJson", funcGetJson)
	p.SetFunc("holdTime", funcHoldTime)
	p.SetFunc("holdCount", funcHoldCount)
	p.SetFunc("setCache", funcSetCache)
	p.SetFunc("getCache", funcGetCache)
	p.SetFunc("speak", funcSpeak)
	p.SetFunc("toInt", funcToInt)
	p.SetFunc("toInt8", funcToInt8)
	p.SetFunc("toInt16", funcToInt16)
	p.SetFunc("toInt32", funcToInt32)
	p.SetFunc("toInt64", funcToInt64)
	p.SetFunc("toUint8", funcToUint8)
	p.SetFunc("toUint16", funcToUint16)
	p.SetFunc("toUint32", funcToUint32)
	p.SetFunc("toUint64", funcToUint64)
	p.SetFunc("toFloat", funcToFloat)
	p.SetFunc("toFloat32", funcToFloat32)
	p.SetFunc("toFloat64", funcToFloat64)
	p.SetFunc("toString", funcToString)
	p.SetFunc("toBool", funcToBool)
	p.SetFunc("toBIN", funcToBIN)
	p.SetFunc("toHEX", funcToHEX)
	p.SetFunc("reboot", funcReboot)
	p.SetFunc("shell", funcShell)
	p.SetFunc("http", funcHTTP)
	p.SetFunc("sum", funcSum)
	return p
}

type Pool struct {
	all  []*Client
	free chan *Client
}

func (this *Pool) take() *Client {
	return <-this.free
}

func (this *Pool) put(c *Client, err ...*error) {
	if e := recover(); e != nil {
		er := fmt.Errorf("%v", e)
		for _, v := range err {
			*v = er
		}
	}
	this.free <- c
}

func (this *Pool) Exec(text string) (_ *conv.Var, err error) {
	c := this.take()
	defer this.put(c, &err)
	return c.Exec(text)
}

func (this *Pool) Set(key string, value interface{}) error {
	for _, v := range this.all {
		if err := v.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (this *Pool) SetFunc(key string, fn func(*Args) interface{}) error {
	for _, v := range this.all {
		if err := v.Set(key, fn); err != nil {
			return err
		}
	}
	return nil
}
