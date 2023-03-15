package sys_save

import (
	"bytes"
	"fmt"
	"github.com/minio/minio-go"
)

func NewMinio(cfg *MinioConfig) (*Minio, error) {
	cli, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, false)
	if err != nil {
		return nil, err
	}
	return &Minio{Client: cli, cfg: cfg}, nil
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
}

func (this *Minio) Save(filename string, fileBytes []byte) (string, error) {
	_, err := this.PutObject(this.cfg.BucketName, filename, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{})
	return fmt.Sprintf("%s/%s/%s", this.cfg.Endpoint, this.cfg.BucketName, filename), err
}
