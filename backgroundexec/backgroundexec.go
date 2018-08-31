// Package backgroundexec is a GopherJS wrapper for allowing execution on iOS when app goes to background (https://github.com/jocull/phonegap-backgroundjs).
//
// Install plugin:
//  cordova plugin add https://github.com/jocull/phonegap-backgroundjs
package backgroundexec

import (
	"github.com/gopherjs/gopherjs/js"
)

// SetBackgroundSeconds allows background execution for seconds
func SetBackgroundSeconds(seconds int) {
	rootObj().Call("setBackgroundSeconds", seconds)
}

// LockBackgroundTime allows background execution until max allowed time (180s)
func LockBackgroundTime() {
	rootObj().Call("lockBackgroundTime")
}

// UnlockBackgroundTime stops background tasks immediately
func UnlockBackgroundTime() {
	rootObj().Call("unlockBackgroundTime")
}

func rootObj() *js.Object {
	return js.Global.Get("plugins").Get("backgroundjs")
}
