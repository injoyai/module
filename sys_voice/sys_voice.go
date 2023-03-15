package sys_voice

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/injoyai/base/g"
	"sync"
)

func Speak(msg string) error {
	return New().Speak(msg)
}

// Save ./wav xxx
func Save(path, msg string) error {
	return New().Save(path, msg)
}

func New() *Voice {
	return &Voice{
		Rate:   0,
		Volume: 100,
	}
}

var mu sync.Mutex

type Voice struct {
	Rate   int //语速
	Volume int //音量
}

func (this *Voice) SetRate(n int) *Voice {
	this.Rate = n
	return this
}

func (this *Voice) SetVolume(n int) *Voice {
	if n > 100 {
		n = 100
	} else if n < 0 {
		n = 0
	}
	this.Volume = n
	return this
}

func (this *Voice) Speak(msg string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	defer g.Recover(&err)
	g.PanicErr(ole.CoInitialize(0))
	unknown, err := oleutil.CreateObject("SAPI.SpVoice")
	g.PanicErr(err)
	voice, err := unknown.QueryInterface(ole.IID_IDispatch)
	g.PanicErr(err)
	_, err = oleutil.PutProperty(voice, "Rate", this.Rate)
	g.PanicErr(err)
	_, err = oleutil.PutProperty(voice, "Volume", this.Volume)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(voice, "Speak", msg)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(voice, "WaitUntilDone", 0)
	g.PanicErr(err)
	voice.Release()
	ole.CoUninitialize()
	return nil
}

func (this *Voice) Save(path, msg string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	defer g.Recover(&err)
	g.PanicErr(ole.CoInitialize(0))
	unknown, err := oleutil.CreateObject("SAPI.SpVoice")
	g.PanicErr(err)
	voice, err := unknown.QueryInterface(ole.IID_IDispatch)
	g.PanicErr(err)
	saveFile, err := oleutil.CreateObject("SAPI.SpFileStream")
	g.PanicErr(err)
	ff, err := saveFile.QueryInterface(ole.IID_IDispatch)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(ff, "Open", path, 3, true)
	g.PanicErr(err)
	_, err = oleutil.PutPropertyRef(voice, "AudioOutputStream", ff)
	g.PanicErr(err)
	_, err = oleutil.PutProperty(voice, "Rate", this.Rate)
	g.PanicErr(err)
	_, err = oleutil.PutProperty(voice, "Volume", this.Volume)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(voice, "Speak", msg)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(voice, "WaitUntilDone", 0)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(ff, "Close")
	g.PanicErr(err)
	ff.Release()
	voice.Release()
	ole.CoUninitialize()
	return nil
}
