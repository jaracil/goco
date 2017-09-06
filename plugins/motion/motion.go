package motion

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type Acceleration struct {
	*js.Object
	X         float64 `js:"x"`
	Y         float64 `js:"y"`
	Z         float64 `js:"z"`
	Timestamp int64   `js:"timestamp"`
}

type Watcher struct {
	id *js.Object
}

func wrapAcceleration(obj *js.Object) *Acceleration {
	return &Acceleration{Object: obj}
}

func CurrentAcceleration() (acc *Acceleration, err error) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		acc = wrapAcceleration(obj)
		close(ch)
	}
	fail := func() {
		err = errors.New("Error getting accelerometer data")
		close(ch)
	}

	js.Global.Get("navigator").Get("accelerometer").Call("getCurrentAcceleration", success, fail)
	<-ch
	return
}

func NewWatcher(cb func(acc *Acceleration, err error), options map[string]interface{}) *Watcher {
	success := func(obj *js.Object) {
		acc := wrapAcceleration(obj)
		cb(acc, nil)
	}

	fail := func() {
		err := errors.New("Error getting accelerometer data")
		cb(nil, err)
	}

	id := js.Global.Get("navigator").Get("accelerometer").Call("watchAcceleration", success, fail, options)
	return &Watcher{id: id}
}

func (w *Watcher) Close() {
	js.Global.Get("navigator").Get("accelerometer").Call("clearWatch", w.id)
}
