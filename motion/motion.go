// Package motion is a GopherJS wrapper for cordova device-motion plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-device-motion
package motion

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

// Acceleration type with x, y, z axes and timestamp
type Acceleration struct {
	*js.Object
	X         float64 `js:"x"`
	Y         float64 `js:"y"`
	Z         float64 `js:"z"`
	Timestamp int64   `js:"timestamp"`
}

// Watcher type monitors acceleration changes
type Watcher struct {
	*js.Object
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("navigator").Get("accelerometer")
	}
	return instance
}

// CurrentAcceleration gets the current acceleration.
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

	mo().Call("getCurrentAcceleration", success, fail)
	<-ch
	return
}

// NewWatcher creates a new motion watcher
func NewWatcher(cb func(acc *Acceleration, err error), options interface{}) *Watcher {
	success := func(a *Acceleration) {
		cb(a, nil)
	}

	fail := func() {
		err := errors.New("Error getting accelerometer data")
		cb(nil, err)
	}

	id := mo().Call("watchAcceleration", success, fail, options)
	return &Watcher{Object: id}
}

// Close cancels watcher
func (w *Watcher) Close() {
	mo().Call("clearWatch", w)
}
