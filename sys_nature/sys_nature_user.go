package sys_nature

import (
	"github.com/injoyai/base/g"
	"time"
)

type User struct {
	ID       int64  `json:"id"`                    //主键
	Name     string `json:"name"`                  //名称
	Memo     string `json:"memo"`                  //备注
	Key      string `json:"key" xorm:"index"`      //全局唯一标识,opc可以根据这个来读取数据
	Type     g.Type `json:"type"`                  //实时值类型,int,string,float,bool
	Value    string `json:"value"`                 //值
	Unit     string `json:"unit"`                  //单位 例℃
	LastTime int64  `json:"lastTime"`              //最后时间,毫秒
	InDate   int64  `json:"inDate" xorm:"created"` //创建时间
}

func (this *User) GetName() string { return this.Name }

func (this *User) GetMemo() string { return this.Memo }

func (this *User) GetValue() interface{} { return this.Value }

func (this *User) GetType() string { return string(this.Type) }

func (this *User) GetLastTime() int64 { return this.LastTime }

func (this *User) Writable() bool { return true }

func (this *User) WriteValue(value string) error {
	this.Value = value
	this.LastTime = time.Now().Unix() / 1e6
	_, err := db.Where("ID", this.ID).Cols("Value,LastTime").Update(this)
	return err
}

type UserSearch struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

func (this *User) GetList(req *UserSearch) ([]*User, int64, error) {

	return nil, 0, nil
}
