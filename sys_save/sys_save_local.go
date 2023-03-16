package sys_save

import (
	"github.com/injoyai/base/bytes/crypt/md5"
	"os"
	"path/filepath"
	"time"
)

func NewLocal(dir string) *Local {
	return &Local{Dir: dir}
}

type Local struct {
	Dir string
}

func (this *Local) Save(filename string, fileBytes []byte, rename ...bool) (string, error) {
	now := time.Now()
	dir := filepath.Join(this.Dir, now.Format("2006-01-02/"))
	if err := os.MkdirAll(dir, 0666); err != nil {
		return "", err
	}
	//判断是否需要重命名
	if len(rename) > 0 && rename[0] {
		filename = md5.Encrypt(filename + now.String())
	}
	filename = filepath.Join(dir, filename)
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(fileBytes)
	return filename, err
}
