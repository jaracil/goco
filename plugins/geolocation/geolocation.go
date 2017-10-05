package geolocation

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco/plugins/cordova"
)

type Coords struct {
	*js.Object
	Latitude         float64 `js:"latitude"`
	Longitude        float64 `js:"longitude"`
	Altitude         float64 `js:"altitude"`
	Accuracy         float64 `js:"accuracy"`
	AltitudeAccuracy float64 `js:"altitudeAccuracy"`
	Heading          float64 `js:"heading"`
	Speed            float64 `js:"speed"`
}

type Position struct {
	*js.Object
	Coords    *Coords `js:"coords"`
	Timestamp int64   `js:"timestamp"`
}

type Watcher struct {
	*js.Object
}

var mo *js.Object

func init() {
	cordova.OnDeviceReady(func() {
		mo = js.Global.Get("navigator").Get("geolocation")
	})
}

func CurrentPosition(options interface{}) (pos *Position, err error) {
	ch := make(chan struct{})
	success := func(p *Position) {
		pos = p
		close(ch)
	}
	fail := func(obj *js.Object) {
		err = errors.New(obj.Get("message").String())
		close(ch)
	}

	mo.Call("getCurrentPosition", success, fail, options)
	<-ch
	return
}

func NewWatcher(cb func(pos *Position, err error), options interface{}) *Watcher {
	success := func(p *Position) {
		cb(p, nil)
	}

	fail := func(obj *js.Object) {
		err := errors.New(obj.Get("message").String())
		cb(nil, err)
	}

	id := mo.Call("watchPosition", success, fail, options)
	return &Watcher{Object: id}
}

func (w *Watcher) Close() {
	mo.Call("clearWatch", w)
}
