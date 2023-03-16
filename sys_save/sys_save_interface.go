package sys_save

type Interface interface {
	Save(filename string, fileBytes []byte, rename ...bool) (string, error)
}
