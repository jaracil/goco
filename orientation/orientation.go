// Package orientation is a GopherJS wrapper for cordova orientation plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-device-orientation
package orientation

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

// Heading contains orientation values.
type Heading struct {
	*js.Object
	MagneticHeading float64 `js:"magneticHeading"` // The heading in degrees from 0-359.99 at a single moment in time.
	TrueHeading     float64 `js:"trueHeading"`     // The heading relative to the geographic North Pole in degrees 0-359.99 at a single moment in time. A negative value indicates that the true heading can't be determined.
	HeadingAccuracy float64 `js:"headingAccuracy"` // The deviation in degrees between the reported heading and the true heading.
	Timestamp       int64   `js:"timestamp"`       // Milliseconds from Unix epoch
}

// Watcher type monitors orientation changes
type Watcher struct {
	*js.Object
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("navigator").Get("compass")
	}
	return instance
}

// CurrentHeading returns current device's heading.
func CurrentHeading() (heading *Heading, err error) {
	ch := make(chan struct{})
	success := func(h *Heading) {
		heading = h
		close(ch)
	}
	fail := func() {
		err = errors.New("Error getting compass data")
		close(ch)
	}

	mo().Call("getCurrentHeading", success, fail)
	<-ch
	return
}

// NewWatcher creates new orientation tracking watcher.
func NewWatcher(cb func(*Heading, error), options map[string]interface{}) *Watcher {
	success := func(h *Heading) {
		cb(h, nil)
	}

	fail := func() {
		err := errors.New("Error getting compass data")
		cb(nil, err)
	}

	id := mo().Call("watchHeading", success, fail, options)
	return &Watcher{Object: id}
}

// Close cancels tracking watcher
func (w *Watcher) Close() {
	mo().Call("clearWatch", w)
}
