package sys_script

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/module/sys_go"
	"github.com/injoyai/module/sys_voice"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	holdTime  = maps.NewSafe()
	holdCount = maps.NewSafe()
	cacheMap  = maps.NewSafe()
)

// funcPrint 打印输出
func funcPrint(debug *[]string) func(args *Args) interface{} {
	return func(args *Args) interface{} {
		msg := fmt.Sprint(args.Interfaces()...)
		*debug = append(*debug, msg)
		return nil
	}
}

// funcPrintf 格式化打印
func funcPrintf(debug *[]string) func(args *Args) interface{} {
	return func(args *Args) interface{} {
		list := args.Args
		if len(list) > 0 {
			msg := fmt.Sprintf(list[0].String(), args.Interfaces()[1:]...)
			*debug = append(*debug, msg)
		}
		return nil
	}
}

/*




 */

// funcGetJson 解析json,读取其中数据
func funcGetJson(args *Args) interface{} {
	return conv.NewMap(args.GetString(1)).GetString(args.GetString(2))
}

// funcSpeak 播放语音
func funcSpeak(args *Args) interface{} {
	msg := args.GetString(1)
	sys_go.Default.Go(&sys_go.Create{
		Group: "",
		Name:  "播放语音",
		Handler: func(ctx context.Context, a *sys_go.Entity, p sys_go.Param) error {
			return sys_voice.Speak(msg)
		},
	})
	return nil
}

// funcBase64ToString base64编码
func funcBase64ToString(args *Args) interface{} {
	data := args.GetString(1)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// funcBase64ToBytes base64解码
func funcBase64ToBytes(args *Args) interface{} {
	data := args.GetString(1)
	bs, _ := base64.StdEncoding.DecodeString(data)
	return string(bs)
}

// funcHEXToBytes 字符转字节 例 "0102" >>> []byte{0x01,0x02}
func funcHEXToBytes(args *Args) interface{} {
	s := args.GetString(1)
	bs, _ := hex.DecodeString(s)
	return string(bs)
}

// funcHEXToString 字节转字符 例 []byte{0x01,0x02} >>> "0102"
func funcHEXToString(args *Args) interface{} {
	return hex.EncodeToString([]byte(args.GetString(1)))
}

// funcHoldTime 连续保持时间触发
func funcHoldTime(args *Args) interface{} {
	key := args.GetString(1)     //key(唯一标识)
	hold := args.GetBool(2)      //保持
	second := args.GetFloat64(3) //持续时间(秒)
	if hold {
		t := time.Now()
		first := holdTime.GetVar(key)
		if !first.IsNil() {
			return first != nil && t.Sub(first.Val().(time.Time)).Seconds() > second
		}
		holdTime.Set(key, t) //第一次触发
		return false
	}
	holdTime.Del(key)
	return false
}

// funcHoldCount 连续保持次数触发
func funcHoldCount(args *Args) interface{} {
	key := args.GetString(1) //key(唯一标识)
	rule := args.GetBool(2)  //规则
	count := args.GetInt(3)  //持续次数
	if rule {
		co := holdCount.GetInt(key)
		co++
		holdCount.Set(key, co)
		return co >= count
	}
	holdCount.Del(key)
	return false
}

// funcSetCache 设置缓存
func funcSetCache(args *Args) interface{} {
	key := args.GetString(1)
	val := args.GetString(2)
	expiration := args.GetInt(3, 0)
	cacheMap.Set(key, val, time.Millisecond*time.Duration(expiration))
	return nil
}

// funcGetCache 获取缓存
func funcGetCache(args *Args) interface{} {
	key := args.GetString(1)
	return cacheMap.MustGet(key)
}

// funcLen 取字符长度
func funcLen(args *Args) interface{} {
	key := args.GetString(1)
	return len(key)
}

// funcToInt 任意类型转int
func funcToInt(args *Args) interface{} {
	return conv.Int(args.GetString(1))
}

// funcToInt8 任意类型转int8
func funcToInt8(args *Args) interface{} {
	return conv.Int8(args.GetString(1))
}

// funcToInt16 任意类型转int16
func funcToInt16(args *Args) interface{} {
	return conv.Int16(args.GetString(1))
}

// funcToInt32 任意类型转int32
func funcToInt32(args *Args) interface{} {
	return conv.Int32(args.GetString(1))
}

// funcToInt64 任意类型转int64
func funcToInt64(args *Args) interface{} {
	return conv.Int64(args.GetString(1))
}

// funcToUint8 任意类型转uint8
func funcToUint8(args *Args) interface{} {
	return conv.Uint8(args.GetString(1))
}

// funcToUint16 任意类型转uint8
func funcToUint16(args *Args) interface{} {
	return conv.Uint16(args.GetString(1))
}

// funcToUint32 任意类型转uint32
func funcToUint32(args *Args) interface{} {
	return conv.Uint32(args.GetString(1))
}

// funcToUint64 任意类型转uint64
func funcToUint64(args *Args) interface{} {
	return conv.Uint32(args.GetString(1))
}

// funcToFloat 任意类型转浮点
func funcToFloat(args *Args) interface{} {
	return conv.Float64(args.GetString(1))
}

// funcToFloat32 任意类型转浮点32位
func funcToFloat32(args *Args) interface{} {
	return conv.Float32(args.GetString(1))
}

// funcToFloat64 任意类型转浮点64位
func funcToFloat64(args *Args) interface{} {
	return conv.Float64(args.GetString(1))
}

// funcToString 任意类型转字符串
func funcToString(args *Args) interface{} {
	return args.GetString(1)
}

// funcToBool 任意类型转bool
func funcToBool(args *Args) interface{} {
	return conv.Bool(args.GetString(1))
}

// funcToBIN 数字转成2进制字符串
func funcToBIN(args *Args) interface{} {
	byte := args.GetInt(2)
	data := interface{}(args.GetInt64(1))
	switch byte {
	case 1:
		data = conv.Uint8(data)
	case 2:
		data = conv.Uint16(data)
	case 4:
		data = conv.Uint32(data)
	case 8:
		data = conv.Uint32(data)
	default:
		data = conv.Uint16(data)
	}
	return conv.BINStr(data)
}

func funcToHEX(args *Args) interface{} {
	data := args.GetInt64(1)
	bytes := []byte(nil)
	switch args.GetInt(2) {
	case 1:
		bytes = []byte{uint8(data)}
	case 2:
		bytes = conv.Bytes(uint16(data))
	case 4:
		bytes = conv.Bytes(uint32(data))
	case 8:
		bytes = conv.Bytes(uint64(data))
	default:
		bytes = conv.Bytes(uint8(data))
	}
	return hex.EncodeToString(bytes)
}

// funcSum 求和
func funcSum(args *Args) interface{} {
	sum := 0
	for _, v := range args.Args {
		sum += v.Int()
	}
	return sum
}

// funcReboot 重启系统
func funcReboot(args *Args) interface{} {
	return shell("reboot")
}

// funcShell 执行脚本
func funcShell(args *Args) interface{} {
	arg := []string(nil)
	for _, v := range args.Args {
		arg = append(arg, v.String())
	}
	return shell(arg...)
}

// funcHTTP http请求,协程执行
func funcHTTP(args *Args) interface{} {
	method := strings.ToUpper(args.GetString(1))
	url := args.GetString(2)
	body := args.GetString(3)
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}

func shell(args ...string) interface{} {
	list := append([]string{"/c"}, args...)
	switch runtime.GOOS {
	case "windows":
		_, err := exec.Command("cmd", list...).CombinedOutput()
		return err
	case "linux":
		list[0] = "-c"
		_, err := exec.Command("bash", list...).CombinedOutput()
		return err
	}
	return errors.New("未知操作系统:" + runtime.GOOS)
}
