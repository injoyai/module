package sys_task

type SysTask struct {
	ID      int64  `json:"id"`                    //
	Group   string `json:"group"`                 //分组
	Name    string `json:"name"`                  //名称
	Memo    string `json:"memo"`                  //备注
	Spec    string `json:"spec"`                  //表达式
	Script  string `json:"script" xorm:"text"`    //脚本
	InDate  int64  `json:"inDate" xorm:"created"` //创建时间
	cornKey string `json:"-" xorm:"-"`            //
}

type SysTaskCreateReq struct {
	Group  string `json:"group"`              //分组
	Name   string `json:"name"`               //名称
	Memo   string `json:"memo"`               //备注
	Spec   string `json:"spec"`               //表达式
	Script string `json:"script" xorm:"text"` //脚本
}

func (this *SysTaskCreateReq) New() (*SysTask, string, error) {
	return &SysTask{
		Group:  this.Group,
		Name:   this.Name,
		Memo:   this.Memo,
		Spec:   this.Spec,
		Script: this.Script,
	}, "Group,Name,Memo,Spec,Script", nil
}

type SysTaskUpdateReq struct {
	ID int64 `json:"id"` //
	*SysTaskCreateReq
}
