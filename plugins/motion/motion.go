package motion

import (
	"errors"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

type Acceleration struct {
	X         float64
	Y         float64
	Z         float64
	Timestamp time.Time
}

type Watcher struct {
	id *js.Object
}

func wrapAcceleration(obj *js.Object) *Acceleration {
	acc := &Acceleration{X: obj.Get("x").Float(), Y: obj.Get("y").Float(), Z: obj.Get("z").Float()}
	uts := obj.Get("timestamp").Int64()
	acc.Timestamp = time.Unix(uts/1000, (uts%1000)*1000000)
	return acc
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
