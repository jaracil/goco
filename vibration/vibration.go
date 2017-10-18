// Package vibration is a GopherJS wrapper for cordova vibration plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-vibration
package vibration

import (
	"github.com/gopherjs/gopherjs/js"
)

// Vibrate has different functionalities based on parameters passed to it.
// For example:
//  Vibrate(3000) // vibrate for 3 seconds
//  Vibrate(3000, 1000, 3000) // vibrate for 3 seconds, wait for 1 second, vibrate for 3 seconds
func Vibrate(pat ...interface{}) {
	js.Global.Get("navigator").Call("vibrate", pat)
}
