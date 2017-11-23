// Package univlinks is a GopherJS wrapper for cordova-universal-links-plugin.
//
// Install plugin:
//  cordova plugin add cordova-universal-links-plugin
//
// (Noncomplete implementation)
package univlinks

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

var mo *js.Object

func getRoot() *js.Object {
	return js.Global.Get("universalLinks")
}

// Subscribe so cb is executed when eventName (defined in config.xml) is emitted when opening a link
func Subscribe(eventName string, cb func(*js.Object)) {
	println(fmt.Sprintf("univlinks.Subscribe: eventName=%q, cb=%v", eventName, cb))
	getRoot().Call("subscribe", eventName, cb)
}
