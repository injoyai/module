package sys_save

import (
	"bytes"
	"fmt"
	"github.com/injoyai/base/bytes/crypt/md5"
	"github.com/minio/minio-go"
	"time"
)

func NewMinio(cfg *MinioConfig) Interface {
	cli, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, false)
	return &Minio{Client: cli, cfg: cfg, err: err}
}

type MinioConfig struct {
	Endpoint   string //地址
	AccessKey  string //访问key
	SecretKey  string //秘钥
	BucketName string //桶名称
}

type Minio struct {
	cfg *MinioConfig
	*minio.Client
	err error
}

func (this *Minio) Save(filename string, fileBytes []byte, rename ...bool) (string, error) {
	now := time.Now()
	if this.err != nil {
		return "", this.err
	}
	if len(rename) > 0 && rename[0] {
		filename = md5.Encrypt(filename + now.String())
	}
	filename = now.Format("2006-01/") + filename
	_, err := this.PutObject(this.cfg.BucketName, filename, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{})
	return fmt.Sprintf("%s/%s/%s", this.cfg.Endpoint, this.cfg.BucketName, filename), err
}
