package background

import (
	"github.com/gopherjs/gopherjs/js"
)

var bgc *js.Object

func bg() *js.Object {
	if bgc == nil {
		bgc = js.Global.Get("cordova").Get("plugins").Get("backgroundMode")
	}
	return bgc
}

func Enable(en bool) {
	bg().Call("setEnabled", en)
}

func InBackground() bool {
	return bg().Call("isActive").Bool()
}

func MoveToBackground() {
	bg().Call("moveToBackground")
}

func MoveToForeground() {
	bg().Call("moveToForeground")
}

func OverrideBackButton() {
	bg().Call("overrideBackButton")
}

func ExcludeFromTaskList() {
	bg().Call("ExcludeFromTaskList")
}

func IsScreenOff() (ret bool) {
	ch := make(chan struct{})
	success := func(b bool) {
		ret = !b
		close(ch)
	}
	bg().Call("isScreenOff", success)
	<-ch
	return
}

func Wakeup() {
	bg().Call("wakeup")
}

func Unlock() {
	bg().Call("unlock")
}

func DisableWebViewOptimizations() {
	bg().Call("disableWebViewOptimizations")
}

func OnEnable(f func()) {
	bg().Call("on", "enable", f)
}

func OnDisable(f func()) {
	bg().Call("on", "disable", f)
}

func OnActivate(f func()) {
	bg().Call("on", "activate", f)
}

func OnDeactivate(f func()) {
	bg().Call("on", "deactivate", f)
}

func OnFailure(f func()) {
	bg().Call("on", "failure", f)
}

func UnEnable(f func()) {
	bg().Call("un", "enable", f)
}

func UnDisable(f func()) {
	bg().Call("un", "disable", f)
}

func UnActivate(f func()) {
	bg().Call("un", "activate", f)
}

func UnDeactivate(f func()) {
	bg().Call("un", "deactivate", f)
}

func UnFailure(f func()) {
	bg().Call("un", "failure", f)
}

func SetDefaults(opts map[string]interface{}) {
	bg().Call("setDefaults", opts)
}

func Configure(opts map[string]interface{}) {
	bg().Call("configure", opts)
}
