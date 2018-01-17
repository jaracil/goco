// Package univlinks is a GopherJS wrapper for cordova-universal-links-plugin.
//
// Install plugin:
//  cordova plugin add cordova-universal-links-plugin
//
// (Noncomplete implementation)
package univlinks

import (
	"github.com/gopherjs/gopherjs/js"
)

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("universalLinks")
	}
	return instance
}

// Subscribe so cb is executed when eventName (defined in config.xml) is emitted when opening a link
func Subscribe(eventName string, cb func(*js.Object)) {
	mo().Call("subscribe", eventName, cb)
}
