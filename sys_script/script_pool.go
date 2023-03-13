package script

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
