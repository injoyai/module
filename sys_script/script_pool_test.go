package sys_script

import (
	"log"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	p := New(10)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				v, err := p.Exec(`
function add(a){
return a+2 
}
print(add(1))
`)
				if err != nil {
					t.Error(err)
				} else {
					t.Log(v.Val())
				}
			}
		}()
	}
	select {}
}

func TestPoolSpend(t *testing.T) {
	p := New(10)
	for i := 0; i < 10000; i++ {
		_, err := p.Exec(`
function add(a){
return a+2 
}
add(1)
`)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestPoolGoSpend(t *testing.T) {
	p := New(1)
	p.SetFunc("sleep", testSleep)
	for i := 0; i < 10000; i++ {
		go func() {
			p.Set("a", "1")
			p.Set("b", 2)
			_, err := p.Exec(`
function add(a){
return a+2 
}
add(1)
sleep(10)
`)
			if err != nil {
				t.Error(err)
			}
		}()
	}
}

func TestClientSpend(t *testing.T) {
	c := NewClient()
	c.SetFunc("sleep", testSleep)
	for i := 0; i < 10000; i++ {
		c.Set("a", "1")
		c.Set("b", 2)
		_, err := c.Exec(`
function add(a){
return a+2 
}
add(1)
sleep(1)
`)
		if err != nil {
			t.Error(err)
		}
	}
}

func testSleep(args *Args) interface{} {
	log.Println(args.GetInt(1))
	time.Sleep(time.Millisecond * time.Duration(args.GetInt(1)))
	return nil
}
