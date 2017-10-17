package vibration

import (
	"github.com/gopherjs/gopherjs/js"
)

func Vibrate(pat ...interface{}) {
	js.Global.Get("navigator").Call("vibrate", pat)
}
