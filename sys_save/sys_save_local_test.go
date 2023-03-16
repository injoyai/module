package sys_save

import (
	"io/ioutil"
	"testing"
)

func TestNewLocal(t *testing.T) {
	c := NewLocal("./")
	bs, err := ioutil.ReadFile("C:\\Users\\injoy\\Pictures\\Camera Roll\\WIN_20220908_08_09_35_Pro.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(c.Save("钱测试", bs))
}
