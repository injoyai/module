package sys_cfg

import (
	"errors"
	"github.com/injoyai/conv"
)

const (
	String = "string"
	Bool   = "bool"
	Int    = "int"
	Float  = "float"
)

// ValueType 基础值类型
type ValueType string

func (this ValueType) String() string {
	switch this {
	case String:
		return "字符"
	case Bool:
		return "布尔"
	case Int:
		return "整数"
	case Float:
		return "浮点"
	}
	return "未知"
}

func (this ValueType) Value(v interface{}) interface{} {
	switch this {
	case String:
		return conv.String(v)
	case Bool:
		return conv.Bool(v)
	case Int:
		return conv.Int(v)
	case Float:
		return conv.Float64(v)
	}
	return v
}

func (this *ValueType) Check() error {
	switch *this {
	case String, Bool, Int, Float:
	case "":
		*this = String
	default:
		return errors.New("未知数据类型")
	}
	return nil
}

type SysCfgSearch struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Name     string `json:"name"`   //配置数据名称
	Key      string `json:"key"`    //配置数据标识
	DeptID   string `json:"deptID"` //部门标识
	Group    string `json:"group"`  //分组标识
}

type SysCfg struct {
	ID     int64 `json:"id"`                    //主键
	InDate int64 `json:"inDate" xorm:"created"` //创建时间

	Name        string      `json:"name" xorm:"index"` //名称
	Memo        string      `json:"memo"`              //备注
	Key         string      `json:"key" xorm:"index"`  //配置数据标识
	ValueString string      `json:"-"`                 //配置数据内容
	ValueType   ValueType   `json:"valueType"`         //配置数据类型
	Value       interface{} `json:"value" xorm:"-"`    //配置数据

	DeptID    string `json:"deptID" xorm:"index"` //部门标识
	Group     string `json:"group" xorm:"index"`  //分组
	GroupName string `json:"groupName"`           //分组名称
	Sort      int    `json:"sort" xorm:"index"`
}

func (this *SysCfg) Resp() *SysCfg {
	if len(this.ValueType) == 0 {
		this.ValueType = String
	}
	this.Value = this.ValueType.Value(this.ValueString)
	return this
}

type SysCfgCreateReq struct {
	Name      string      `json:"name"`             //名称
	Memo      string      `json:"memo"`             //备注
	Key       string      `json:"key" xorm:"index"` //配置数据标识
	ValueType ValueType   `json:"valueType"`        //配置数据类型
	Value     interface{} `json:"value" xorm:"-"`   //配置数据

	DeptID    string `json:"deptID" xorm:"index"` //部门id
	Group     string `json:"group" xorm:"index"`  //分组
	GroupName string `json:"groupName"`           //分组名称
	Sort      int    `json:"sort" xorm:"index"`
}

func (this *SysCfgCreateReq) New() (*SysCfg, error) {
	if err := this.ValueType.Check(); err != nil {
		return nil, err
	}
	return &SysCfg{
		Name:        this.Name,
		Memo:        this.Memo,
		Key:         this.Key,
		ValueString: conv.String(this.Value),
		ValueType:   this.ValueType,

		DeptID:    this.DeptID,
		Group:     this.Group,
		GroupName: this.GroupName,
		Sort:      this.Sort,
	}, nil
}

type SysCfgUpdateReq struct {
	ID int64 `json:"id"`
	*SysCfgCreateReq
}

func (this *SysCfg) Update(req *SysCfgUpdateReq) (string, error) {
	data, err := req.New()
	if err != nil {
		return "", err
	}
	this.Name = data.Name
	this.Memo = data.Memo
	this.Key = data.Key
	this.ValueString = data.ValueString
	this.ValueType = data.ValueType
	this.Group = data.Group
	this.GroupName = data.GroupName
	this.Sort = data.Sort
	return "Name,Memo,Key,ValueString,ValueType,Group,GroupName,Sort", nil
}
