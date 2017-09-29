package background

import (
	"github.com/gopherjs/gopherjs/js"
)

func Enable(en bool) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("setEnabled", en)
}

func InBackground() bool {
	return js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("isActive").Bool()
}

func MoveToBackground() {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("moveToBackground")
}

func MoveToForeground() {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("moveToForeground")
}

func OverrideBackButton() {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("overrideBackButton")
}

func ExcludeFromTaskList() {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("ExcludeFromTaskList")
}

func IsScreenOff() (ret bool) {
	ch := make(chan struct{})
	success := func(b bool) {
		ret = !b
		close(ch)
	}
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("isScreenOff", success)
	<-ch
	return
}

func Wakeup() {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("wakeup")
}

func Unlock() {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("unlock")
}

func DisableWebViewOptimizations() {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("disableWebViewOptimizations")
}

func OnEnable(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("on", "enable", f)
}

func OnDisable(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("on", "disable", f)
}

func OnActivate(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("on", "activate", f)
}

func OnDeactivate(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("on", "deactivate", f)
}

func OnFailure(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("on", "failure", f)
}

func UnEnable(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("un", "enable", f)
}

func UnDisable(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("un", "disable", f)
}

func UnActivate(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("un", "activate", f)
}

func UnDeactivate(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("un", "deactivate", f)
}

func UnFailure(f func()) {
	js.Global.Get("cordova").Get("plugins").Get("backgroundMode").Call("un", "failure", f)
}
