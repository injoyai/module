package sys_cfg

import (
	"errors"
	"github.com/injoyai/base/maps"
	"xorm.io/xorm"
)

/*
New
系统配置模块,缓存
根据单位配置,
根据分组配置,
根据key配置
*/
func New(db *xorm.Engine) (*Entity, error) {
	e := &Entity{db: db, Safe: maps.NewSafe()}
	return e, e.loading()
}

// Entity 系统配置实例
type Entity struct {
	db *xorm.Engine
	*maps.Safe
}

// loading 数据加载,初始化
func (this *Entity) loading() error {
	data, err := this.GetAll()
	if err != nil {
		return err
	}
	for _, v := range data {
		this.Set(v.ID, v)
	}
	return nil
}

// GetAll 获取全部配置
func (this *Entity) GetAll() (data []*SysCfg, err error) {
	err = this.db.Find(&data)
	return data, err
}

// GetList 获取配置列表
func (this *Entity) GetList(req *SysCfgSearch) ([]*SysCfg, int64, error) {
	data := []*SysCfg{}
	session := this.db.Desc("Sort")
	if req.PageSize > 0 {
		session.Limit(req.PageSize, req.PageNum*req.PageSize)
	}
	if len(req.DeptID) > 0 {
		session.Where("DeptID=?", req.DeptID)
	}
	if len(req.Group) > 0 {
		session.Where("Group=?", req.Group)
	}
	if len(req.Key) > 0 {
		session.Where("`Key` like ?", "%"+req.Key+"%")
	}
	if len(req.Name) > 0 {
		session.Where("Name like ?", "%"+req.Name+"%")
	}
	co, err := session.FindAndCount(&data)
	return data, co, err
}

// GetByCache 获取缓存的数据,不存在则尝试从数据库获取
func (this *Entity) GetByCache(id int64) (*SysCfg, error) {
	v, err := this.GetOrSetByHandler(id, func() (interface{}, error) {
		//缓存不存在,尝试去数据库查询
		return this.GetByID(id)
	})
	if err != nil {
		return nil, err
	}
	return v.(*SysCfg), nil
}

// GetByID 根据主键获取配置数据
func (this *Entity) GetByID(id int64) (*SysCfg, error) {
	data := new(SysCfg)
	has, err := this.db.Where("ID=?", id).Get(data)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("配置数据不存在")
	}
	return data, nil
}

// GetByKey 根据key获取配置数据
func (this *Entity) GetByKey(key string) (*SysCfg, error) {
	data := new(SysCfg)
	has, err := this.db.Where("`Key`=?", key).Get(data)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("配置数据不存在")
	}
	return data, nil
}

// Exist 判断配置数据标识是否存在
func (this *Entity) Exist(deptID, key string, query interface{}, args ...interface{}) (bool, error) {
	return this.db.Where("deptID=? and `Key`=?", deptID, key).
		Where(query, args...).Exist(new(SysCfg))
}

// Post 新建配置
func (this *Entity) Post(req *SysCfgCreateReq) error {
	//数据整理,校验
	data, err := req.New()
	if err != nil {
		return err
	}
	//判断数据库配置标识是否存在
	if has, err := this.Exist(req.DeptID, req.Key, nil); err != nil {
		return err
	} else if has {
		return errors.New("配置数据标识已存在")
	}
	//添加到数据库
	_, err = this.db.Insert(data)
	if err == nil {
		//添加到缓存
		this.Set(data.ID, data)
	}
	return err
}

// Put 修改配置
func (this *Entity) Put(req *SysCfgUpdateReq) error {
	//获取数据库数据
	data, err := this.GetByCache(req.ID)
	if err != nil {
		return err
	}
	//数据整理,校验
	cols, err := data.Update(req)
	if err != nil {
		return err
	}
	//判断数据库配置标识是否存在
	if has, err := this.Exist(req.DeptID, req.Key, "ID<>?", req.ID); err != nil {
		return err
	} else if has {
		return errors.New("配置数据标识已存在")
	}
	_, err = this.db.Where("ID=?", req.ID).Cols(cols).Update(data)
	if err == nil {
		//更新到缓存
		this.Set(data.ID, data)
	}
	return err
}

// PutList 批量修改配置
func (this *Entity) PutList(req []*SysCfgUpdateReq) error {
	for _, v := range req {
		if err := this.Put(v); err != nil {
			return err
		}
	}
	return nil
}

// DelByID 根据主键删除配置
func (this *Entity) DelByID(id int64) error {
	//删除数据库
	_, err := this.db.Where("ID=?", id).Delete(new(SysCfg))
	if err == nil {
		//删除缓存
		this.Del(id)
	}
	return err
}

// DelByDeptID 删除部门的全部配置数据
func (this *Entity) DelByDeptID(deptID string) error {
	_, err := this.db.Where("DeptID=?", deptID).Delete(new(SysCfg))
	if err == nil {
		this.Range(func(key, value interface{}) bool {
			if value.(*SysCfg).DeptID == deptID {
				this.Del(key)
			}
			return true
		})
	}
	return err
}

// DelByKey 删除部门的一个key
func (this *Entity) DelByKey(deptID, key string) error {
	_, err := this.db.Where("DeptID=? and `Key`=?", deptID, key).Delete(new(SysCfg))
	if err == nil {
		this.Range(func(key, value interface{}) bool {
			x := value.(*SysCfg)
			if x.DeptID == deptID && x.Key == key {
				this.Del(key)
			}
			return true
		})
	}
	return err
}

// DelByGroup 删除部门的一个配置数据分组
func (this *Entity) DelByGroup(deptID, group string) error {
	_, err := this.db.Where("DeptID=? and Group=?", deptID, group).Delete(new(SysCfg))
	if err == nil {
		this.Range(func(key, value interface{}) bool {
			x := value.(*SysCfg)
			if x.DeptID == deptID && x.Group == group {
				this.Del(key)
			}
			return true
		})
	}
	return err
}
