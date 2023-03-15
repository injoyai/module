package sys_moddle

import (
	"bytes"
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/cache"
	"github.com/injoyai/conv"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// NewMonitorAPI 监控api实例
func NewMonitorAPI(cfg *cache.FileLogConfig, i cache.IFileLog) *MonitorAPI {
	return &MonitorAPI{FileLog: cache.NewFileLog(cfg), i: i}
}

type MonitorAPI struct {
	*cache.FileLog
	i cache.IFileLog
}

func (this *MonitorAPI) Middle(r *http.Request) func() {
	//过滤Options请求
	if r.Method == http.MethodOptions {
		return func() {}
	}
	//websocket 不参与耗时计算
	if r.Header.Get("Upgrade") == "websocket" {
		this.WriteAny(NewMonitorAPILog(r, time.Now()).Json())
		return func() {}
	}
	t := time.Now()
	return func() {
		this.WriteAny(NewMonitorAPILog(r, t).Json())
	}
}

// GetAPIInfoList 获取接口请求列表
func (this *MonitorAPI) GetAPIInfoList(pageSize int) (interface{}, error) {
	list, err := this.GetLogLast(pageSize)
	if err != nil {
		return nil, err
	}
	result := []*MonitorAPILog{}
	for _, v := range list {
		result = append(result, DecodeMonitorAPILog(v))
	}
	return result, nil
}

func (this *MonitorAPI) GetAPIRate(start, end int64) (interface{}, error) {
	list, err := this.GetAPILog(start, end)
	if err != nil {
		return nil, err
	}
	m := map[string][]*MonitorAPILog{}
	for _, v := range list {
		key := fmt.Sprintf("%s#%s", v.Method, v.URI)
		m[key] = append(m[key], v)
	}
	result := []*MonitorAPIRate(nil)
	for key, v := range m {
		l := strings.Split(key, "#")
		r := &MonitorAPIRate{
			Method: l[0],
			URI:    l[1],
		}
		for _, k := range v {
			r.Add(k)
		}
		result = append(result, r)
	}
	sort.Sort(MonitorAPISort(result))
	return result, nil
}

// GetAPICurve 获取曲线
func (this *MonitorAPI) GetAPICurve(start, end int64) (_ interface{}, err error) {
	return this.GetLogCurve(time.Unix(start, 0), time.Unix(end, 0), time.Hour, this.i)
}

// GetAPILog 获取记录
func (this *MonitorAPI) GetAPILog(start, end int64) ([]*MonitorAPILog, error) {
	result := []*MonitorAPILog{}
	list, err := this.GetLog(time.Unix(start, 0), time.Unix(end, 0))
	for _, v := range list {
		result = append(result, DecodeMonitorAPILog(v))
	}
	return result, err
}

func NewMonitorAPILog(r *http.Request, start time.Time) *MonitorAPILog {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	if len(body) > 200 {
		body = append(body[:200], []byte("...")...)
	}
	now := time.Now()
	return &MonitorAPILog{
		Method: r.Method,
		URI:    r.URL.Path,
		Query:  r.URL.RawQuery,
		Body:   string(body),
		Spend:  int64(now.Sub(start)) / 1e6,
		InDate: now.Unix(),
	}
}

func DecodeMonitorAPILog(bs []byte) *MonitorAPILog {
	s := string(bs)
	if len(s) > 0 {
		s = s[1 : len(s)-1]
	}
	a := new(MonitorAPILog)
	for _, v := range strings.Split(s, ",") {
		if list := strings.SplitN(v, ":", 2); len(list) == 2 {
			val := list[1]
			switch list[0] {
			case `"m"`:
				if len(val) >= 2 {
					a.Method = val[1 : len(val)-1]
				}
			case `"u"`:
				if len(val) >= 2 {
					a.URI = val[1 : len(val)-1]
				}
			case `"q"`:
				if len(val) >= 2 {
					a.Query = val[1 : len(val)-1]
				}
			case `"b"`:
				if len(val) >= 2 {
					a.Body = val[1 : len(val)-1]
				}
			case `"s"`:
				n, _ := strconv.Atoi(val)
				a.Spend = int64(n)
			case `"at"`:
				n, _ := strconv.Atoi(val)
				a.InDate = int64(n)
			}
		}
	}
	return a
}

// MonitorAPIRate 接口请求频率信息
type MonitorAPIRate struct {
	Method   string `json:"method"`   //请求方式
	URI      string `json:"path"`     //路径
	Spend    int64  `json:"spend"`    //平均耗时
	MaxSpend int64  `json:"maxSpend"` //最大耗时
	Number   int64  `json:"number"`   //请求的次数
}

func (this *MonitorAPIRate) Add(a *MonitorAPILog) {
	this.Spend = (this.Spend*(this.Number) + a.Spend) / (this.Number + 1)
	atomic.AddInt64(&this.Number, 1)
	if a.Spend > this.MaxSpend {
		this.MaxSpend = a.Spend
	}
}

type MonitorAPILog struct {
	Method string `json:"method"` //请求方式
	URI    string `json:"path"`   //路径
	Query  string `json:"query"`  //query
	Body   string `json:"body"`   //body
	Spend  int64  `json:"spend"`  //耗时
	InDate int64  `json:"inDate"` //创建时间
}

func (this *MonitorAPILog) GetSecond() int64 {
	return this.InDate
}

func (this *MonitorAPILog) Json() string {
	return conv.String(g.Map{"m": this.Method, "u": this.URI, "q": this.Query, "b": this.Body, "s": this.Spend, "at": this.InDate})
}

func (this *MonitorAPILog) NodeHour() int64 {
	return this.InDate - this.InDate%(60*60)
}

type MonitorAPISort []*MonitorAPIRate

func (this MonitorAPISort) Less(i, j int) bool {
	return this[i].Number > this[j].Number
}

func (this MonitorAPISort) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this MonitorAPISort) Len() int {
	return len(this)
}
