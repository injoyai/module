package sys_script

import (
	"fmt"
	"github.com/injoyai/conv"
	"github.com/robertkrimen/otto"
	"sync"
)

func NewClient() *Client {
	vm := otto.New()
	cli := &Client{
		Otto: vm,
	}
	cli.Set("print", cli.toFunc(func(args *Args) interface{} {
		fmt.Println(args.Interfaces()...)
		return nil
	}))
	cli.Exec("var console={\nlog:function(any){\nprint(any)\n}\n}")
	return cli
}

// Client 万次执行0.11s
type Client struct {
	*otto.Otto
	mu sync.Mutex
}

func (this *Client) Exec(text string) (*conv.Var, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	value, err := this.Otto.Run(text)
	if err != nil {
		return conv.Nil(), err
	}
	val, _ := value.Export()
	return conv.New(val), nil
}

func (this *Client) GetVar(key string) *conv.Var {
	val, _ := this.Otto.Get(key)
	value, _ := val.Export()
	return conv.New(value)
}

func (this *Client) Set(key string, value interface{}) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	switch fn := value.(type) {
	case func(*Args) interface{}:
		value = this.toFunc(fn)
	case func():
		value = this.toFunc(func(args *Args) interface{} {
			fn()
			return nil
		})
	}
	return this.Otto.Set(key, value)
}

func (this *Client) SetFunc(key string, fn func(*Args) interface{}) error {
	return this.Set(key, this.toFunc(fn))
}

func (this *Client) toFunc(fn func(*Args) interface{}) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		it, _ := call.This.Export()
		args := []*conv.Var(nil)
		for _, v := range call.ArgumentList {
			val, _ := v.Export()
			args = append(args, conv.New(val))
		}
		arg := &Args{
			This:   conv.New(it),
			Args:   args,
			Client: this,
		}
		result, _ := otto.ToValue(fn(arg))
		return result
	}
}
