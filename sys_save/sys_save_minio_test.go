package sys_save

import (
	"io/ioutil"
	"testing"
)

func TestMinio_Save(t *testing.T) {

	c, err := NewMinio(&MinioConfig{
		Endpoint:   "192.168.10.103:9002",
		AccessKey:  "minioadmin",
		SecretKey:  "minioadmin",
		BucketName: "qianlang-iot-admin",
	})
	if err != nil {
		t.Error(err)
		return
	}
	bs, err := ioutil.ReadFile("C:\\Users\\injoy\\Pictures\\Camera Roll\\WIN_20220908_08_09_35_Pro.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(c.Save("钱测试", bs))
}
