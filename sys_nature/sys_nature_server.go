package sys_nature

var entity = NewSysNature()

func GetNature(group, key string) *Nature {
	n, ok := entity.Get(group, key)
	if !ok {
		return NewNatureNull(group, key)
	}
	return NewNature(group, n)
}

func SetNature(group, key string, value INature) {
	entity.Set(group, key, value)
}
