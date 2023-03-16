package sys_save

import (
	"bytes"
	"github.com/injoyai/base/bytes/crypt/md5"
	"github.com/injoyai/base/g"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"time"
)

func NewQiniu(cfg *QiniuConfig) Interface {
	return &Qiniu{cfg: cfg}
}

type QiniuConfig struct {
	AccessKey string //访问key
	SecretKey string //秘钥
	Domain    string //前缀
	Space     string //空间
}

type Qiniu struct {
	cfg *QiniuConfig
}

// Save 上传数据
// @name,名称
// @r,读数据流
func (this *Qiniu) Save(filename string, fileBytes []byte, rename ...bool) (string, error) {
	if len(rename) > 0 && rename[0] {
		filename = md5.Encrypt(filename + time.Now().String())
	}
	buff := bytes.NewBuffer(fileBytes)
	mac := qbox.NewMac(this.cfg.AccessKey, this.cfg.SecretKey)
	putPolicy := storage.PutPolicy{Scope: this.cfg.Space}
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS:      false,                // 是否使用https域名
		UseCdnDomains: false,                // 上传是否使用CDN上传加速
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.Put(g.Ctx(), &ret, upToken, filename, buff, int64(len(fileBytes)), &storage.PutExtra{})
	return this.cfg.Domain + ret.Key, err
}
