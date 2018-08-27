// Package geolocation is a GopherJS wrapper for cordova geolocation plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-geolocation
package geolocation

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

// Coords defines a set of geographic coordinates.
type Coords struct {
	*js.Object
	Latitude         float64 `js:"latitude"`         // Latitude in decimal degrees.
	Longitude        float64 `js:"longitude"`        // Longitude in decimal degrees.
	Altitude         float64 `js:"altitude"`         // Height of the position in meters above the ellipsoid.
	Accuracy         float64 `js:"accuracy"`         // Accuracy level of the latitude and longitude coordinates in meters.
	AltitudeAccuracy float64 `js:"altitudeAccuracy"` // Accuracy level of the altitude coordinate in meters.
	Heading          float64 `js:"heading"`          // Direction of travel, specified in degrees counting clockwise relative to the true north.
	Speed            float64 `js:"speed"`            // Current ground speed of the device, specified in meters per second.
}

// Position defines coordinates and timestamp.
type Position struct {
	*js.Object
	Coords    *Coords `js:"coords"`
	Timestamp int64   `js:"timestamp"` // Milliseconds from Unix epoch
}

// Watcher type monitors position changes
type Watcher struct {
	*js.Object
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("navigator").Get("geolocation")
	}
	return instance
}

// CurrentPosition returns current device's position.
// Options must be map[string]interface{} type.
// See https://cordova.apache.org/docs/en/latest/reference/cordova-plugin-geolocation/index.html#options for available options.
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

	mo().Call("getCurrentPosition", success, fail, options)
	<-ch
	return
}

// NewWatcher creates new position tracking watcher.
// Options must be map[string]interface{} type.
// See https://cordova.apache.org/docs/en/latest/reference/cordova-plugin-geolocation/index.html#options for available options.
func NewWatcher(cb func(pos *Position, err error), options interface{}) *Watcher {
	success := func(p *Position) {
		cb(p, nil)
	}

	fail := func(obj *js.Object) {
		err := errors.New(obj.Get("message").String())
		cb(nil, err)
	}

	id := mo().Call("watchPosition", success, fail, options)
	return &Watcher{Object: id}
}

// Close cancels tracking watcher
func (w *Watcher) Close() {
	mo().Call("clearWatch", w)
}
