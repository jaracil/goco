package motion

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco/plugins/cordova"
)

type Acceleration struct {
	*js.Object
	X         float64 `js:"x"`
	Y         float64 `js:"y"`
	Z         float64 `js:"z"`
	Timestamp int64   `js:"timestamp"`
}

type Watcher struct {
	*js.Object
}

var mo *js.Object

func init() {
	cordova.OnDeviceReady(func() {
		mo = js.Global.Get("navigator").Get("accelerometer")
	})
}

func CurrentAcceleration() (acc *Acceleration, err error) {
	ch := make(chan struct{})
	success := func(a *Acceleration) {
		acc = a
		close(ch)
	}
	fail := func() {
		err = errors.New("Error getting accelerometer data")
		close(ch)
	}

	mo.Call("getCurrentAcceleration", success, fail)
	<-ch
	return
}

func NewWatcher(cb func(acc *Acceleration, err error), options map[string]interface{}) *Watcher {
	success := func(a *Acceleration) {
		cb(a, nil)
	}

	fail := func() {
		err := errors.New("Error getting accelerometer data")
		cb(nil, err)
	}

	id := mo.Call("watchAcceleration", success, fail, options)
	return &Watcher{Object: id}
}

func (w *Watcher) Close() {
	mo.Call("clearWatch", w)
}
